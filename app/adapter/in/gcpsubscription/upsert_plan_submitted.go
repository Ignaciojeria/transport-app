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
		newPlanSubmitted,
		subscriptionwrapper.NewSubscriptionManager,
		usecase.NewUpsertPlan)
}
func newPlanSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	upsertPlan usecase.UpsertPlan,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-upsert-plan-submitted"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 5
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		var input request.UpsertPlanRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}

		inputMapped := input.Map()
		//inputMapped.Organization.SetKey(m.Attributes["organization"])

		_, err := upsertPlan(ctx, inputMapped)
		if err != nil {
			m.Ack()
			return http.StatusOK, err
		}
		m.Ack()
		return http.StatusOK, err
	}
	go sm.WithMessageProcessor(messageProcessor).
		Start(subscriptionRef)
	return messageProcessor
}
