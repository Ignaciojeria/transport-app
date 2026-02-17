package subscriber

import (
	"context"
	"errors"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"

	"micartapro/app/adapter/out/imagegenerator"
	"micartapro/app/adapter/out/storage"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/creem"
	"micartapro/app/usecase/menu"
)

func init() {
	ioc.Registry(
		newPullAllEvents,
		eventprocessing.NewSubscriberStrategy,
		observability.NewObservability,
		menu.NewOnMenuInteractionRequest,
		menu.NewOnMenuCreateRequest,
		menu.NewOnUserMenusInsertedWebhook,
		menu.NewOnImageGenerationRequest,
		menu.NewOnImageEditionRequest,
		creem.NewOnCreemSubscriptionTrialingWebhook,
		creem.NewOnCreemCheckoutCompletedWebhook,
		creem.NewOnCreemSubscriptionActiveWebhook,
		creem.NewOnCreemSubscriptionPaidWebhook,
		creem.NewOnCreemSubscriptionCanceledWebhook,
		creem.NewOnCreemSubscriptionExpiredWebhook,
		creem.NewOnCreemSubscriptionUpdateWebhook,
		creem.NewOnCreemSubscriptionPausedWebhook,
		creem.NewOnCreemRefundCreatedWebhook,
		creem.NewOnCreemDisputeCreatedWebhook,
	)
}

// Pull all events from the event bus
func newPullAllEvents(
	sub eventprocessing.Subscriber,
	obs observability.Observability,
	onMenuInteractionRequest menu.OnMenuInteractionRequest,
	onMenuCreateRequest menu.OnMenuCreateRequest,
	onUserMenusInsertedWebhook menu.OnUserMenusInsertedWebhook,
	onImageGenerationRequest menu.OnImageGenerationRequest,
	onImageEditionRequest menu.OnImageEditionRequest,
	onCreemSubscriptionTrialingWebhook creem.OnCreemSubscriptionTrialingWebhook,
	onCreemCheckoutCompletedWebhook creem.OnCreemCheckoutCompletedWebhook,
	onCreemSubscriptionActiveWebhook creem.OnCreemSubscriptionActiveWebhook,
	onCreemSubscriptionPaidWebhook creem.OnCreemSubscriptionPaidWebhook,
	onCreemSubscriptionCanceledWebhook creem.OnCreemSubscriptionCanceledWebhook,
	onCreemSubscriptionExpiredWebhook creem.OnCreemSubscriptionExpiredWebhook,
	onCreemSubscriptionUpdateWebhook creem.OnCreemSubscriptionUpdateWebhook,
	onCreemSubscriptionPausedWebhook creem.OnCreemSubscriptionPausedWebhook,
	onCreemRefundCreatedWebhook creem.OnCreemRefundCreatedWebhook,
	onCreemDisputeCreatedWebhook creem.OnCreemDisputeCreatedWebhook) eventprocessing.MessageProcessor {
	subscriptionName := "micartapro.events.v3"
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

		case events.EventCreemCheckoutCompletedWebhook:
			var request events.CreemCheckoutCompletedWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onCreemCheckoutCompletedWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_creem_checkout_completed_webhook",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventCreemSubscriptionActiveWebhook:
			var request events.CreemSubscriptionActiveWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onCreemSubscriptionActiveWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_creem_subscription_active_webhook",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventCreemSubscriptionPaidWebhook:
			var request events.CreemSubscriptionPaidWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onCreemSubscriptionPaidWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_creem_subscription_paid_webhook",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventCreemSubscriptionCanceledWebhook:
			var request events.CreemSubscriptionCanceledWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onCreemSubscriptionCanceledWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_creem_subscription_canceled_webhook",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventCreemSubscriptionExpiredWebhook:
			var request events.CreemSubscriptionExpiredWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onCreemSubscriptionExpiredWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_creem_subscription_expired_webhook",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventCreemSubscriptionUpdateWebhook:
			var request events.CreemSubscriptionUpdateWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onCreemSubscriptionUpdateWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_creem_subscription_update_webhook",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventCreemSubscriptionPausedWebhook:
			var request events.CreemSubscriptionPausedWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onCreemSubscriptionPausedWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_creem_subscription_paused_webhook",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventCreemRefundCreatedWebhook:
			var request events.CreemRefundCreatedWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onCreemRefundCreatedWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_creem_refund_created_webhook",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventCreemDisputeCreatedWebhook:
			var request events.CreemDisputeCreatedWebhook
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
			}
			err := onCreemDisputeCreatedWebhook(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_creem_dispute_created_webhook",
					"error", err.Error(),
				)
				return http.StatusInternalServerError
			}

		case events.EventImageGenerationRequested:
			var request events.ImageGenerationRequestEvent
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
				return http.StatusBadRequest
			}
			err := onImageGenerationRequest(spanCtx, request)
			if err != nil {
				obs.Logger.Error("error_processing_image_generation_request",
					"error", err.Error(),
					"menuId", request.MenuID,
					"menuItemId", request.MenuItemID,
				)
				// ACK para no reintentar (evita retries infinitos en fallos de API/upload)
				return http.StatusAccepted
			}

		case events.EventImageEditionRequested:
			var request events.ImageEditionRequestEvent
			if err := event.DataAs(&request); err != nil {
				obs.Logger.Error("failed_to_deserialize_cloudevent",
					"error", err.Error(),
				)
				return http.StatusBadRequest
			}
			err := onImageEditionRequest(spanCtx, request)
			if err != nil {
				// Imagen no disponible o vacía: ACK para no reintentar (no es candidata a edición)
				if errors.Is(err, imagegenerator.ErrReferenceImageNotAvailable) {
					obs.Logger.WarnContext(spanCtx, "reference_image_not_available_ack",
						"error", err.Error(),
						"menuId", request.MenuID,
						"menuItemId", request.MenuItemID,
					)
					return http.StatusAccepted
				}
				obs.Logger.Error("error_processing_image_edition_request",
					"error", err.Error(),
					"menuId", request.MenuID,
					"menuItemId", request.MenuItemID,
				)
				return http.StatusInternalServerError
			}
		}
		return http.StatusOK
	}

	// Start subscriber
	go sub.Start(subscriptionName, processor, eventprocessing.ReceiveSettings{MaxOutstandingMessages: 3})
	return processor
}
