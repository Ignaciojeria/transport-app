package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/domain"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"

	canonicaljson "transport-app/app/shared/caonincaljson"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newRegistrationSubmittedConsumer,
		natsconn.NewJetStream,
		usecase.NewCreateAccountWorkflow,
		usecase.NewCreateTenantWorkflow,
		usecase.NewAssociateTenantAccountWorkflow,
		usecase.NewCreateDefaultClientCredentialsWorkflow,
		usecase.NewSendClientCredentialsEmailWorkflow,
		observability.NewObservability,
		configuration.NewConf,
	)
}

func newRegistrationSubmittedConsumer(
	js jetstream.JetStream,
	createAccountWorkflow usecase.CreateAccountWorkflow,
	createTenantWorkflow usecase.CreateTenantWorkflow,
	associateTenantAccountWorkflow usecase.AssociateTenantAccountWorkflow,
	createDefaultClientCredentialsWorkflow usecase.CreateDefaultClientCredentialsWorkflow,
	sendClientCredentialsEmailWorkflow usecase.SendClientCredentialsEmailWorkflow,
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
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.REGISTRATION_SUBMITTED_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.REGISTRATION_SUBMITTED_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.registrationSubmitted",
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

		// Deserializar el payload como RegisterRequest
		var input request.RegisterRequest
		if err := json.Unmarshal(pubsubMsg.Data, &input); err != nil {
			obs.Logger.Error("Error deserializando payload de registro", "error", err)
			msg.Ack()
			return
		}

		// Generar idempotency key único basado en el email para todo el proceso de registro
		registrationKey, err := canonicaljson.HashKey(ctx, "registration", input.Email)
		if err != nil {
			obs.Logger.ErrorContext(ctx, "Error generando key para registro", "error", err)
			msg.Nak()
			return
		}
		registrationCtx := sharedcontext.WithIdempotencyKey(ctx, registrationKey)

		// Procesar el registro con idempotency key único
		account := domain.Account{
			Email: input.Email,
		}
		if err := createAccountWorkflow(registrationCtx, account); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando registro", "error", err)
			msg.Nak() // Usar Nak para permitir reintentos
			return
		}

		tenant := domain.Tenant{
			ID:      account.UUID(),
			Name:    input.OrganizationName,
			Country: countries.ByName(input.Country),
		}
		if err := createTenantWorkflow(registrationCtx, tenant); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando registro", "error", err)
			msg.Nak() // Usar Nak para permitir reintentos
			return
		}

		tenantAccount := domain.TenantAccount{
			Tenant:  tenant,
			Account: account,
			Role:    "owner",
		}
		if err := associateTenantAccountWorkflow(registrationCtx, tenantAccount); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando registro", "error", err)
			msg.Nak() // Usar Nak para permitir reintentos
			return
		}

		clientCredentials := domain.ClientCredentials{
			ID:            uuid.New(),
			TenantID:      tenant.ID,
			TenantCountry: tenant.Country,
			ClientID:      account.UUID().String(),
			ClientSecret:  uuid.New().String(),
			AllowedScopes: []string{
				"orders:read",
				"orders:write",
				"routes:read",
				"routes:write",
				"nodes:read",
				"nodes:write",
				"optimization:read",
				"optimization:write",
				"deliveries:read",
				"deliveries:write",
			},
		}
		if err := createDefaultClientCredentialsWorkflow(registrationCtx, clientCredentials); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando registro", "error", err)
			msg.Nak() // Usar Nak para permitir reintentos
			return
		}

		if err := sendClientCredentialsEmailWorkflow(registrationCtx, clientCredentials.ClientID, account.Email); err != nil {
			obs.Logger.ErrorContext(ctx, "Error procesando registro", "error", err)
			msg.Nak() // Usar Nak para permitir reintentos
			return
		}

		obs.Logger.InfoContext(ctx, "Registro procesado exitosamente desde NATS",
			"eventType", "registrationSubmitted",
			"email", input.Email)
		msg.Ack()
	})
}
