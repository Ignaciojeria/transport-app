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
		newVehicleSubmitted,
		subscriptionwrapper.NewSubscriptionManager,
		usecase.NewUpsertVehicle)
}
func newVehicleSubmitted(
	sm subscriptionwrapper.SubscriptionManager,
	upsert usecase.UpsertVehicle,
) subscriptionwrapper.MessageProcessor {
	subscriptionName := "transport-app-events-vehicle-submitted"
	subscriptionRef := sm.Subscription(subscriptionName)
	subscriptionRef.ReceiveSettings.MaxOutstandingMessages = 1
	messageProcessor := func(ctx context.Context, m *pubsub.Message) (int, error) {
		var input request.UpsertVehicleRequest
		if err := json.Unmarshal(m.Data, &input); err != nil {
			m.Ack()
			return http.StatusAccepted, err
		}
		domainOBJ := input.Map()
		// Completar datos adicionales desde los atributos del mensaje
		domainOBJ.Commerce = m.Attributes["commerce"]
		domainOBJ.Consumer = m.Attributes["consumer"]
		//domainOBJ.Organization.SetKey(m.Attributes["organization"])
		if err := upsert(ctx, domainOBJ); err != nil {
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
