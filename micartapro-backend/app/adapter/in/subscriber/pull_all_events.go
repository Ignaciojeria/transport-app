package subscriber

import (
	"context"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	"micartapro/app/domain"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/usecase"
)

func init() {
	ioc.Registry(
		newPullAllEvents,
		eventprocessing.NewSubscriberStrategy,
		observability.NewObservability,
		usecase.NewMenuInteraction)
}

// Pull all events from the event bus
func newPullAllEvents(
	sub eventprocessing.Subscriber,
	obs observability.Observability,
	menuInteraction usecase.MenuInteraction) eventprocessing.MessageProcessor {
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

		switch event.Type() {
		case domain.EventMenuInteractionRequested:
			var request domain.MenuInteractionRequest
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			_, err := menuInteraction(ctx, request)
			if err != nil {
				obs.Logger.Error("error_processing_menu_interaction",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}
		}
		return http.StatusOK
	}

	// Start subscriber
	go sub.Start(subscriptionName, processor)
	return processor
}
