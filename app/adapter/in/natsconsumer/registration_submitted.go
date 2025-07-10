package natsconsumer

import (
	"context"
	"encoding/json"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/usecase"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newRegistrationSubmittedConsumer,
		natsconn.NewJetStream,
		usecase.NewRegister,
		observability.NewObservability,
		configuration.NewConf,
	)
}

func newRegistrationSubmittedConsumer(
	js jetstream.JetStream,
	register usecase.Register,
	obs observability.Observability,
	conf configuration.Conf,
) (jetstream.ConsumeContext, error) {
	// Validación para verificar si el nombre de la suscripción está vacío
	if conf.REGISTRATION_SUBMITTED_SUBSCRIPTION == "" {
		obs.Logger.Warn("Registration submitted subscription name is empty, skipping consumer initialization")
		// Retornar nil para indicar que no hay consumidor activo
		return nil, nil
	}

	ctx := context.Background()
	consumer, err := js.CreateOrUpdateConsumer(ctx, conf.TRANSPORT_APP_TOPIC, jetstream.ConsumerConfig{
		Name:          conf.REGISTRATION_SUBMITTED_SUBSCRIPTION,
		Durable:       conf.REGISTRATION_SUBMITTED_SUBSCRIPTION,
		FilterSubject: conf.TRANSPORT_APP_TOPIC + ".*.*.registrationSubmitted",
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

		// Deserializar el payload como RegisterRequest
		var input request.RegisterRequest
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de registro", "error", err)
			msg.Ack()
			return
		}

		// Procesar el registro
		if err := register(ctx, domain.TenantAccount{
			Tenant: domain.Tenant{
				Country: countries.ByName(input.Country),
			},
			Account: domain.Account{
				Email: input.Email,
			},
			Role: "owner",
		}); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando registro", "error", err)
			msg.Ack()
			return
		}

		obs.Logger.InfoContext(ctx, "Registro procesado exitosamente desde NATS",
			"eventType", "registrationSubmitted",
			"email", input.Email)
		msg.Ack()
	})
}
