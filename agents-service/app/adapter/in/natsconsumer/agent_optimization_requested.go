package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
)

func init() {
	ioc.Registry(
		newAgentOptimizationRequested,
		natsconn.NewJetStream,
		observability.NewObservability,
		configuration.NewConf,
		storjbucket.NewTransportAppBucket,
	)
}

func newAgentOptimizationRequested(
	js jetstream.JetStream,
	obs observability.Observability,
	conf configuration.Conf,
	storjBucket *storjbucket.TransportAppBucket,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.AGENT_OPTIMIZATION_REQUESTED_SUBSCRIPTION == "" {
		obs.Logger.Warn("Webhook submitted subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.WEBHOOK_SUBMITTED_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.WEBHOOK_SUBMITTED_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.agentOptimizationRequested",
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
		var input interface{}
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de webhook", "error", err)
			msg.Ack()
			return
		}

		obs.Logger.InfoContext(ctx, "Agent optimization request received", "input", input)
		msg.Ack()
	})
}
