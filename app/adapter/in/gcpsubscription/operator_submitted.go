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
		operator.Organization.OrganizationCountryID = organizationCountryID

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
