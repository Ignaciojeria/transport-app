package natsconsumer

import (
	"context"
	"encoding/json"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/usecase"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newTenantSubmittedConsumer,
		natsconn.NewJetStream,
		usecase.NewCreateTenantAccount,
		observability.NewObservability,
		configuration.NewConf,
	)
}

func newTenantSubmittedConsumer(
	js jetstream.JetStream,
	createTenantAccount usecase.CreateTenantAccount,
	obs observability.Observability,
	conf configuration.Conf,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.TENANT_SUBMITTED_SUBSCRIPTION == "" {
		obs.Logger.Warn("Tenant submitted subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          conf.TENANT_SUBMITTED_SUBSCRIPTION,
		Durable:       conf.TENANT_SUBMITTED_SUBSCRIPTION,
		FilterSubject: conf.TRANSPORT_APP_TOPIC + ".*.*.tenantSubmitted",
		MaxAckPending: 5,
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

		// Deserializar el payload como CreateTenantRequest
		var input request.CreateTenantRequest
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de tenant", "error", err)
			msg.Ack()
			return
		}

		// Procesar el tenant
		if err := createTenantAccount(ctx, input.Map()); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando tenant", "error", err)
			msg.Ack()
			return
		}

		obs.Logger.InfoContext(ctx, "Tenant procesado exitosamente desde NATS",
			"eventType", "tenantSubmitted")
		msg.Ack()
	})
}
