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
		usecase.NewMenuInteraction,
		usecase.NewMenuInteractionCreate)
}

// Pull all events from the event bus
func newPullAllEvents(
	sub eventprocessing.Subscriber,
	obs observability.Observability,
	menuInteraction usecase.MenuInteraction,
	menuInteractionCreate usecase.MenuInteractionCreate) eventprocessing.MessageProcessor {
	subscriptionName := "micartapro.events-sub"
	processor := func(ctx context.Context, event cloudevents.Event) int {

		spanCtx, span := obs.Tracer.Start(contextFromCloudEvent(ctx, event), "pull_all_events")
		defer span.End()

		var input interface{}
		if err := event.DataAs(&input); err != nil {
			obs.Logger.Error("failed_to_unmarshal_cloudevent",
				"error", err.Error(),
			)
			// Invalid payload → ACK (no retry)
			return http.StatusAccepted
		}

		obs.Logger.InfoContext(spanCtx, "unmarshalled_cloudevent",
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
			_, err := menuInteraction(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_menu_interaction",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}
		case domain.EventMenuCreateRequested:
			var request domain.MenuCreateRequest
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := menuInteractionCreate(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_menu_create",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}
		}
		return http.StatusOK
	}

	// Start subscriber
	go sub.Start(subscriptionName, processor, eventprocessing.ReceiveSettings{MaxOutstandingMessages: 1})
	return processor
}
