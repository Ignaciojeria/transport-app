package natsconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase"
	template "transport-app/app/usecase/repository_template"

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
		template.NewCreateRepositoryWorkflow,
		observability.NewObservability,
		configuration.NewConf,
	)
}

func newTenantSubmittedConsumer(
	js jetstream.JetStream,
	createTenantAccount usecase.CreateTenantAccount,
	createRepositoryWorkflow template.CreateRepositoryWorkflow,
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
		Name:          fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.TENANT_SUBMITTED_SUBSCRIPTION),
		Durable:       fmt.Sprintf("%s-%s", conf.ENVIRONMENT, conf.TENANT_SUBMITTED_SUBSCRIPTION),
		FilterSubject: conf.TRANSPORT_APP_TOPIC + "." + conf.ENVIRONMENT + ".*.*.tenantSubmitted",
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

		// Crear repositorio asociado al tenant
		tenantData := input.Map()
		repoName := fmt.Sprintf("tenant-%s", tenantData.Tenant.ID)

		// Generar idempotency key para el repositorio
		repoKey := fmt.Sprintf("create-repository-%s", repoName)
		repoCtx := sharedcontext.WithIdempotencyKey(ctx, repoKey)

		if err := createRepositoryWorkflow(repoCtx, repoName); err != nil {
			obs.Logger.ErrorContext(ctx, "Error creando repositorio para tenant",
				"error", err,
				"tenant_id", tenantData.Tenant.ID,
				"repository_name", repoName)
			// No hacemos return aquí para no fallar el procesamiento del tenant
		} else {
			obs.Logger.InfoContext(ctx, "Repositorio creado exitosamente para tenant",
				"tenant_id", tenantData.Tenant.ID,
				"repository_name", repoName)
		}

		msg.Ack()
	})
}
