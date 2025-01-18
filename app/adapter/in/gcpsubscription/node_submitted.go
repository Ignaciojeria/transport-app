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
)

func init() {
	ioc.Registry(
		newNodeSubmitted,
		subscriptionwrapper.NewSubscriptionManager,
		usecase.NewUpsertNode)
}
func newNodeSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	upsert usecase.UpsertNode,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-node-submitted"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 5
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {

		var input request.UpsertNodeRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}

		if err := upsert(ctx, input.Map()); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}

		m.Ack()
		return http.StatusOK, nil
	}
	go sm.WithMessageProcessor(messageProcessor).
		WithPushHandler("/subscription/" + subscriptionName).
		Start(subscriptionRef)
	return messageProcessor
}
