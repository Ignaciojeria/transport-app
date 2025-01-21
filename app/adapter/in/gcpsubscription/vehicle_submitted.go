package gcpsubscription

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		newVehicleSubmitted,
		subscriptionwrapper.NewSubscriptionManager)
}
func newVehicleSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-vehicle-submitted"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 1
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		var input interface{}
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		fmt.Println("works")
		m.Ack()
		return http.StatusOK, nil
	}
	go sm.WithMessageProcessor(messageProcessor).
		WithPushHandler("/subscription/" + subscriptionName).
		Start(subscriptionRef)
	return messageProcessor
}
