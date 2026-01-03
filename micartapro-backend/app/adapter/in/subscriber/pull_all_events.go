package subscriber

import (
	"context"
	"errors"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	"micartapro/app/adapter/out/storage"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase"
)

func init() {
	ioc.Registry(
		newPullAllEvents,
		eventprocessing.NewSubscriberStrategy,
		observability.NewObservability,
		usecase.NewOnMenuInteractionRequest,
		usecase.NewOnMenuCreateRequest,
		usecase.NewOnUserMenusInsertedWebhook,
		usecase.NewOnCreemSubscriptionTrialingWebhook,
	)
}

// Pull all events from the event bus
func newPullAllEvents(
	sub eventprocessing.Subscriber,
	obs observability.Observability,
	onMenuInteractionRequest usecase.OnMenuInteractionRequest,
	onMenuCreateRequest usecase.OnMenuCreateRequest,
	onUserMenusInsertedWebhook usecase.OnUserMenusInsertedWebhook,
	onCreemSubscriptionTrialingWebhook usecase.OnCreemSubscriptionTrialingWebhook) eventprocessing.MessageProcessor {
	subscriptionName := "micartapro.events.v2"
	processor := func(ctx context.Context, event cloudevents.Event) int {

		ctx = sharedcontext.ContextFromCloudEvent(ctx, event)

		spanCtx, span := obs.Tracer.Start(ctx, "pull_all_events")
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
		case events.EventMenuInteractionRequested:
			var request events.MenuInteractionRequest
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			_, err := onMenuInteractionRequest(spanCtx, request)
			if err == nil {
				break
			}

			// Si el error es "menu_not_found", hacer ACK (no reintentar)
			if errors.Is(err, storage.ErrMenuNotFound) {
				obs.Logger.WarnContext(spanCtx, "menu_not_found_ack",
					"error", err.Error(),
					"menuID", request.MenuID,
				)
				return http.StatusAccepted
			}

			obs.Logger.Error("error_processing_menu_interaction",
				"error", err.Error(),
			)
			return http.StatusInternalServerError
		case events.EventMenuCreateRequested:
			var request events.MenuCreateRequest
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onMenuCreateRequest(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_menu_create",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventUserMenusInsertedWebhook:
			var request events.UserMenusInsertedWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onUserMenusInsertedWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_user_menus_inserted_webhook",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventCreemSubscriptionTrialingWebhook:
			var request events.CreemSubscriptionTrialingWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onCreemSubscriptionTrialingWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_creem_subscription_trialing_webhook",
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
