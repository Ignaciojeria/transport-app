package gcpsubscription

import (
	"context"
	"encoding/json"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/usecase/workers"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newOptimizationRequested,
		subscriptionwrapper.NewSubscriptionManager,
		configuration.NewConf,
		workers.NewFleetOptimizer,
		observability.NewObservability,
	)
}
func newOptimizationRequested(
	sm subscriptionwrapper.SubscriptionManager,
	conf configuration.Conf,
	optimize workers.FleetOptimizer,
	obs observability.Observability,
) subscriptionwrapper.MessageProcessor {
	return func(ctx context.Context, m *pubsub.Message) (int, error) {
		return http.StatusAccepted, nil
	}
	subscriptionName := conf.OPTIMIZATION_REQUESTED_SUBSCRIPTION

	// Validación para verificar si el nombre de la suscripción está vacío
	if subscriptionName == "" {
		obs.Logger.Warn("Optimization requested subscription name is empty, skipping message processor initialization")
		return func(ctx context.Context, m *pubsub.Message) (int, error) {
			return http.StatusAccepted, nil
		}
	}

	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 15
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		if m.Attributes["eventType"] != "optimizationRequested" {
			m.Ack()
			return http.StatusAccepted, nil
		}
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(m.Attributes))
		var input request.OptimizeFleetRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}

		err := optimize(ctx, input.Map())

		if err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}

		//observability.Logger.InfoContext(ctx, "Optimization requested", "res", res)
		m.Ack()
		return http.StatusOK, nil
	}
	go sm.WithMessageProcessor(messageProcessor).
		Start(subscriptionRef)
	return messageProcessor
}
