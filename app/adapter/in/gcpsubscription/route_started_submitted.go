package gcpsubscription

import (
	"context"
	"encoding/json"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"
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
		configuration.NewConf)
}

func newRouteStartedSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	routeStarted usecase.RouteStarted,
	conf configuration.Conf,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := conf.ROUTE_STARTED_SUBMITTED_SUBSCRIPTION
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
