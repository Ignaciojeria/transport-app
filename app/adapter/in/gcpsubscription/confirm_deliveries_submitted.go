package gcpsubscription

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		newConfirmDeliveriesSubmitted,
		subscriptionwrapper.NewSubscriptionManager)
}
func newConfirmDeliveriesSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-confirm-deliveries-submitted"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 5
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		var input request.ConfirmDeliveriesRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		fmt.Println("works")
		m.Ack()
		return http.StatusOK, nil
	}
	go sm.WithMessageProcessor(messageProcessor).
		Start(subscriptionRef)
	return messageProcessor
}
