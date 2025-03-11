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
		newOperatorSubmitted,
		subscriptionwrapper.NewSubscriptionManager,
		usecase.NewCreateAccountOperator)
}
func newOperatorSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	upsertOperator usecase.CreateAccountOperator,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-operator-submitted"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 5
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		var input request.CreateAccountOperatorRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		operator := input.Map()
		operator.Organization.SetKey(m.Attributes["organization"])
		if _, err := upsertOperator(ctx, operator); err != nil {
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
