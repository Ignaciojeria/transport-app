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
		newConfirmDeliveriesSubmitted,
		subscriptionwrapper.NewSubscriptionManager,
		usecase.NewConfirmDeliveries)
}
func newConfirmDeliveriesSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	confirmDeliveries usecase.ConfirmDeliverires,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-confirm-deliveries-submitted"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 5
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(m.Attributes))
		var input request.ConfirmDeliveriesRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		if err := confirmDeliveries(ctx, input.Map(m.Data)); err != nil {
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
