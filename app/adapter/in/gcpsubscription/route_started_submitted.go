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
		newRouteStartedSubmitted,
		subscriptionwrapper.NewSubscriptionManager,
		usecase.NewRouteStarted,
		configuration.NewConf,
		observability.NewObservability)
}

func newRouteStartedSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	routeStarted usecase.RouteStarted,
	conf configuration.Conf,
	obs observability.Observability,
) subscriptionwrapper.MessageProcessor {
	return func(ctx context.Context, m *pubsub.Message) (int, error) {
		return http.StatusAccepted, nil
	}
	subscriptionName := conf.ROUTE_STARTED_SUBMITTED_SUBSCRIPTION

	// Validación para verificar si el nombre de la suscripción está vacío
	if subscriptionName == "" {
		obs.Logger.Warn("Route started submitted subscription name is empty, skipping message processor initialization")
		return func(ctx context.Context, m *pubsub.Message) (int, error) {
			return http.StatusAccepted, nil
		}
	}

	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 10
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		if m.Attributes["eventType"] != "routeStarted" {
			m.Ack()
			return http.StatusAccepted, nil
		}
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(m.Attributes))
		var input request.RouteStartedRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		err := routeStarted(ctx, input.Map(ctx))
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
