package subscriber

import (
	"context"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/observability"
)

func init() {
	ioc.Registry(
		newPullAllEvents,
		eventprocessing.NewSubscriberStrategy,
		observability.NewObservability)
}

// Pull all events from the event bus
func newPullAllEvents(
	sub eventprocessing.Subscriber,
	obs observability.Observability) eventprocessing.MessageProcessor {
	subscriptionName := "micartapro.events-sub"

	processor := func(ctx context.Context, event cloudevents.Event) int {

		var input interface{}
		if err := event.DataAs(&input); err != nil {
			obs.Logger.Error("failed_to_unmarshal_cloudevent",
				"error", err.Error(),
			)
			// Invalid payload → ACK (no retry)
			return http.StatusAccepted
		}
		obs.Logger.Info("unmarshalled_cloudevent",
			"event", event,
		)

		// TODO: call your use case
		// processMenuInteraction(input)

		return http.StatusOK
	}

	// Start subscriber
	go sub.Start(subscriptionName, processor)
	return processor
}
