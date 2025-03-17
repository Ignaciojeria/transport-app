package gcpsubscription

import (
	"context"
	"encoding/json"
	"net/http"
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		newOrderDeliveryConfirmationLog,
		subscriptionwrapper.NewSubscriptionManager)
}
func newOrderDeliveryConfirmationLog(
	sm subscriptionwrapper.SubscriptionManager,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-order-delivery-confirmation-log"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 5
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		var input interface{}
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		m.Ack()
		return http.StatusOK, nil
	}
	go sm.WithMessageProcessor(messageProcessor).
		Start(subscriptionRef)
	return messageProcessor
}
