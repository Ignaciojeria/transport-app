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
		newOrderSubmitted,
		subscriptionwrapper.NewSubscriptionManager,
		usecase.NewCreateOrder,
	)
}
func newOrderSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	createOrder usecase.CreateOrder,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-order-submitted"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 1
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		var input request.UpsertOrderRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		// Mapear al dominio desde el request
		domainOBJ := input.Map()

		// Completar datos adicionales desde los atributos del mensaje
		domainOBJ.Commerce = m.Attributes["commerce"]
		domainOBJ.Consumer = m.Attributes["consumer"]

		orgCountryID, err := strconv.ParseInt(m.Attributes["organizationCountryID"], 10, 64)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("invalid organizationCountryID: %w", err)
		}
		domainOBJ.Organization.OrganizationCountryID = orgCountryID
		if _, err := createOrder(ctx, domainOBJ); err != nil {
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
