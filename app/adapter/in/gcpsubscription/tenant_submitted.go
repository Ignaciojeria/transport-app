package gcpsubscription

import (
	"context"
	"encoding/json"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/usecase"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newTenantSubmitted,
		configuration.NewConf,
		usecase.NewCreateTenantAccount,
		subscriptionwrapper.NewSubscriptionManager,
		observability.NewObservability)
}
func newTenantSubmitted(
	conf configuration.Conf,
	createTenantAccount usecase.CreateTenantAccount,
	sm subscriptionwrapper.SubscriptionManager,
	obs observability.Observability,
) subscriptionwrapper.MessageProcessor {
	return func(ctx context.Context, m *pubsub.Message) (int, error) {
		return http.StatusAccepted, nil
	}
	subscriptionName := conf.TENANT_SUBMITTED_SUBSCRIPTION

	// Validación para verificar si el nombre de la suscripción está vacío
	if subscriptionName == "" {
		obs.Logger.Warn("Tenant submitted subscription name is empty, skipping message processor initialization")
		return func(ctx context.Context, m *pubsub.Message) (int, error) {
			return http.StatusAccepted, nil
		}
	}

	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 10
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		// Filtro defensivo para evitar procesamiento incorrecto
		if m.Attributes["eventType"] != "tenantSubmitted" {
			m.Ack()
			return http.StatusAccepted, nil
		}
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(m.Attributes))
		var input request.CreateTenantRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		err := createTenantAccount(ctx, input.Map())
		if err != nil {
			m.Nack()
			return http.StatusAccepted, err
		}
		m.Ack()
		return http.StatusOK, nil
	}
	go sm.WithMessageProcessor(messageProcessor).
		Start(subscriptionRef)
	return messageProcessor
}
