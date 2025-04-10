package gcpsubscription

import (
	"context"
	"encoding/json"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"
	"transport-app/app/usecase"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func init() {
	ioc.Registry(
		newOrderSubmitted,
		subscriptionwrapper.NewSubscriptionManager,
		usecase.NewCreateOrder,
	)
}
func newOrderSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	createOrder usecase.CreateOrder,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-order-submitted"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 1
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(m.Attributes))
		var input request.UpsertOrderRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		if err := createOrder(ctx, input.Map()); err != nil {
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
