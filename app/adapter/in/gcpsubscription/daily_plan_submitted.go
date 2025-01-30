package gcpsubscription

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/gcppubsub/subscriptionwrapper"
	"transport-app/app/usecase"

	"cloud.google.com/go/pubsub"
	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		newDailyPlanSubmitted,
		subscriptionwrapper.NewSubscriptionManager,
		usecase.NewUpsertDailyPlan)
}
func newDailyPlanSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	upsertDailyPlan usecase.UpsertDailyPlan,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-daily-plan-submitted"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 5
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		var input request.CreateDailyPlanRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}

		organizationCountryIDStr, ok := m.Attributes["organizationCountryID"]
		if !ok {
			m.Ack()
			return http.StatusAccepted, fmt.Errorf("organizationCountryID not found in attributes")
		}

		organizationCountryID, err := strconv.ParseInt(organizationCountryIDStr, 10, 64)
		if err != nil {
			m.Ack()
			return http.StatusAccepted, fmt.Errorf("invalid organizationCountryID: %v", err)
		}
		inputMapped := input.Map()
		inputMapped.Organization.OrganizationCountryID = organizationCountryID

		_, err = upsertDailyPlan(ctx, inputMapped)
		if err != nil {
			m.Ack()
			return http.StatusOK, err
		}
		m.Ack()
		return http.StatusOK, err
	}
	go sm.WithMessageProcessor(messageProcessor).
		WithPushHandler("/subscription/" + subscriptionName).
		Start(subscriptionRef)
	return messageProcessor
}
