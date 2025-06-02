package gcpsubscription

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		newDeliveriesSubmitted,
		subscriptionwrapper.NewSubscriptionManager,
		configuration.NewConf)
}
func newDeliveriesSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	conf configuration.Conf,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := conf.DELIVERIES_SUBMITTED_SUBSCRIPTION
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 5
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		if m.Attributes["eventType"] != "deliveriesSubmitted" {
			m.Ack()
			return http.StatusAccepted, nil
		}
		var input interface{}
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		m.Ack()
		fmt.Println("works")
		return http.StatusOK, nil
	}
	go sm.WithMessageProcessor(messageProcessor).
		WithPushHandler("/subscription/" + subscriptionName).
		Start(subscriptionRef)
	return messageProcessor
}
