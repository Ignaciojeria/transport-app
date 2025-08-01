package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"webhooks/app/adapter/in/fuegoapi/model"
	"webhooks/app/client"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newFleetOptimizedWebhookSubmitted,
		natsconn.NewJetStream,
		observability.NewObservability,
		configuration.NewConf,
		client.NewPostWebhook,
		natsconn.NewKeyValue,
	)
}

func newFleetOptimizedWebhookSubmitted(
	js jetstream.JetStream,
	obs observability.Observability,
	conf configuration.Conf,
	postWebhook client.PostWebhook,
	kv jetstream.KeyValue,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.FLEET_OPTIMIZED_WEBHOOK_SUBSCRIPTION == "" {
		obs.Logger.Warn("Webhook submitted subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.WEBHOOK_SUBMITTED_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.WEBHOOK_SUBMITTED_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.fleetOptimizedWebhook",
		MaxAckPending: 5,
		// Configuración de reintentos: 3 reintentos con intervalos de 2 segundos
		MaxDeliver: 4, // 1 intento inicial + 3 reintentos = 4 total
		BackOff:    []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second},
	})

	if err != nil {
		return nil, err
	}

	return consumer.Consume(func(msg jetstream.Msg) {
		// Deserializar el mensaje como pubsub.Message
		var pubsubMsg pubsub.Message
		if err := json.Unmarshal(msg.Data(), &pubsubMsg); err != nil {
			obs.Logger.Error("Error deserializando mensaje NATS", "error", err)
			msg.Ack()
			return
		}

		// Deserializar el payload como UpsertWebhookRequest
		var input model.FleetOptimizedWebhookBody
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de webhook", "error", err)
			msg.Ack()
			return
		}

		wh := domain.Webhook{
			Type: "fleet-optimized",
		}

		ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(pubsubMsg.Attributes))

		bytes, err := kv.Get(ctx, wh.DocID(ctx).String())
		if err != nil {
			obs.Logger.Error("Error obteniendo webhook", "error", err)
			msg.Ack()
			return
		}

		var webhook domain.Webhook
		if err := json.Unmarshal(bytes.Value(), &webhook); err != nil {
			obs.Logger.Error("Error deserializando webhook", "error", err)
			msg.Ack()
			return
		}
		accessToken := msg.Headers().Get("X-Access-Token")
		webhook.Body = input
		webhook.Headers["X-Access-Token"] = accessToken
		webhook.Headers["tenant"] = msg.Headers().Get("tenant")

		if err := postWebhook(ctx, webhook); err != nil {
			obs.Logger.Error("Error posteando webhook", "error", err)
			msg.Ack()
			return
		}

		obs.Logger.InfoContext(ctx, "Webhook procesado exitosamente desde NATS")
		msg.Ack()
	})
}
