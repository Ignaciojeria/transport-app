package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"

	canonicaljson "transport-app/app/shared/caonincaljson"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newWebhookSubmittedConsumer,
		natsconn.NewJetStream,
		usecase.NewUpsertWebhookWorkflow,
		observability.NewObservability,
		configuration.NewConf,
	)
}

func newWebhookSubmittedConsumer(
	js jetstream.JetStream,
	upsertWebhookWorkflow usecase.UpsertWebhookWorkflow,
	obs observability.Observability,
	conf configuration.Conf,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.WEBHOOK_SUBMITTED_SUBSCRIPTION == "" {
		obs.Logger.Warn("Webhook submitted subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.WEBHOOK_SUBMITTED_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.WEBHOOK_SUBMITTED_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.webhookSubmitted",
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

		// Extraer contexto de OpenTelemetry
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(pubsubMsg.Attributes))

		// Deserializar el payload como UpsertWebhookRequest
		var input request.UpsertWebhookRequest
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de webhook", "error", err)
			msg.Ack()
			return
		}

		// Generar idempotency key único basado en el tipo de webhook
		webhookKey, err := canonicaljson.HashKey(ctx, "webhook", input.Type)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para webhook", "error", err)
			msg.Nak()
			return
		}
		webhookCtx := sharedcontext.WithIdempotencyKey(ctx, webhookKey)

		// Procesar el webhook
		if err := upsertWebhookWorkflow(webhookCtx, input.Map(ctx)); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando webhook", "error", err)
			msg.Nak()
			return
		}

		obs.Logger.InfoContext(ctx, "Webhook procesado exitosamente desde NATS",
			"eventType", "webhookSubmitted",
			"type", input.Type)
		msg.Ack()
	})
}
