package fuegoapi

import (
	"encoding/json"
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"net/http"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(creemSubscriptionWebhook, httpserver.New,
		observability.NewObservability,
		eventprocessing.NewPublisherStrategy,
		apimiddleware.NewValidateCreemWebhookSecretMiddleware,
	)
}

type CreemSubscriptionWebhookPayload struct {
	ID        string          `json:"id"`
	EventType string          `json:"eventType"`
	CreatedAt int64           `json:"created_at"`
	Object    json.RawMessage `json:"object"`
}

// extractUserIDFromMetadata busca el user_id en diferentes lugares posibles de la metadata
func extractUserIDFromMetadata(bodyBytes []byte) string {
	var obj map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		return ""
	}

	// Buscar en object.metadata.user_id
	if object, ok := obj["object"].(map[string]interface{}); ok {
		if metadata, ok := object["metadata"].(map[string]interface{}); ok {
			if userID, ok := metadata["user_id"].(string); ok && userID != "" {
				return userID
			}
		}
		// Buscar en object.subscription.metadata.user_id
		if subscription, ok := object["subscription"].(map[string]interface{}); ok {
			if metadata, ok := subscription["metadata"].(map[string]interface{}); ok {
				if userID, ok := metadata["user_id"].(string); ok && userID != "" {
					return userID
				}
			}
		}
		// Buscar en object.checkout.metadata.user_id
		if checkout, ok := object["checkout"].(map[string]interface{}); ok {
			if metadata, ok := checkout["metadata"].(map[string]interface{}); ok {
				if userID, ok := metadata["user_id"].(string); ok && userID != "" {
					return userID
				}
			}
		}
	}

	return "763a590a-9b8e-4a91-b8ee-47f2a64d003d" //default user id
}

func creemSubscriptionWebhook(
	s httpserver.Server,
	obs observability.Observability,
	publisherManager eventprocessing.PublisherManager,
	validateCreemWebhookSecretMiddleware apimiddleware.ValidateCreemWebhookSecretMiddleware) {
	fuego.Post(s.Manager, "/webhooks/creem",
		func(c fuego.ContextWithBody[any]) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "creemWebhook")
			defer span.End()

			requestBodyAny, err := c.Body()
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			// Convertir a JSON para extraer eventType
			bodyBytes, err := json.Marshal(requestBodyAny)
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error marshaling request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			var payload CreemSubscriptionWebhookPayload
			if err := json.Unmarshal(bodyBytes, &payload); err != nil {
				return nil, fuego.HTTPError{
					Title:  "error unmarshaling request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			obs.Logger.Info("creemWebhook request received",
				"eventID", payload.ID,
				"eventType", payload.EventType,
			)

			// Extraer userID de la metadata una sola vez
			userID := extractUserIDFromMetadata(bodyBytes)
			spanCtx = sharedcontext.WithIdempotencyKey(spanCtx, payload.ID)
			if userID != "" {
				spanCtx = sharedcontext.WithUserID(spanCtx, userID)
			}

			// Identificar el tipo de evento y convertir al tipo correcto
			switch payload.EventType {
			case "checkout.completed":
				var webhook events.CreemCheckoutCompletedWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling checkout.completed webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}
				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.checkout.completed.webhook",
					OrderingKey: webhook.Object.Subscription.ID,
					Event:       webhook,
				}); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error publishing event",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				obs.Logger.InfoContext(spanCtx, "creemCheckoutCompletedWebhook event published", "webhook", webhook)
				return http.StatusOK, nil

			case "subscription.active":
				var webhook events.CreemSubscriptionActiveWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling subscription.active webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}
				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.subscription.active.webhook",
					OrderingKey: webhook.Object.ID,
					Event:       webhook,
				}); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error publishing event",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				obs.Logger.InfoContext(spanCtx, "creemSubscriptionActiveWebhook event published", "webhook", webhook)
				return http.StatusOK, nil

			case "subscription.paid":
				var webhook events.CreemSubscriptionPaidWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling subscription.paid webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}
				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.subscription.paid.webhook",
					OrderingKey: webhook.Object.ID,
					Event:       webhook,
				}); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error publishing event",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				obs.Logger.InfoContext(spanCtx, "creemSubscriptionPaidWebhook event published", "webhook", webhook)
				return http.StatusOK, nil

			case "subscription.canceled":
				var webhook events.CreemSubscriptionCanceledWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling subscription.canceled webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}
				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.subscription.canceled.webhook",
					OrderingKey: webhook.Object.ID,
					Event:       webhook,
				}); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error publishing event",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				obs.Logger.InfoContext(spanCtx, "creemSubscriptionCanceledWebhook event published", "webhook", webhook)
				return http.StatusOK, nil

			case "subscription.expired":
				var webhook events.CreemSubscriptionExpiredWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling subscription.expired webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}
				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.subscription.expired.webhook",
					OrderingKey: webhook.Object.ID,
					Event:       webhook,
				}); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error publishing event",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				obs.Logger.InfoContext(spanCtx, "creemSubscriptionExpiredWebhook event published", "webhook", webhook)
				return http.StatusOK, nil

			case "subscription.update":
				var webhook events.CreemSubscriptionUpdateWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling subscription.update webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}
				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.subscription.update.webhook",
					OrderingKey: webhook.Object.ID,
					Event:       webhook,
				}); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error publishing event",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				obs.Logger.InfoContext(spanCtx, "creemSubscriptionUpdateWebhook event published", "webhook", webhook)
				return http.StatusOK, nil

			case "subscription.trialing":
				var webhook events.CreemSubscriptionTrialingWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling subscription.trialing webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}
				if userID == "" {
					obs.Logger.Error("user_id_missing_in_webhook_metadata", "eventID", webhook.ID)
					return nil, fuego.HTTPError{
						Title:  "invalid webhook payload",
						Detail: "user_id missing in metadata",
						Status: http.StatusBadRequest,
					}
				}
				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.subscription.trialing.webhook",
					OrderingKey: webhook.Object.ID,
					Event:       webhook,
				}); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error publishing event",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				obs.Logger.InfoContext(spanCtx, "creemSubscriptionTrialingWebhook event published", "webhook", webhook)
				return http.StatusOK, nil

			case "subscription.paused":
				var webhook events.CreemSubscriptionPausedWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling subscription.paused webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}
				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.subscription.paused.webhook",
					OrderingKey: webhook.Object.ID,
					Event:       webhook,
				}); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error publishing event",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				obs.Logger.InfoContext(spanCtx, "creemSubscriptionPausedWebhook event published", "webhook", webhook)
				return http.StatusOK, nil

			case "refund.created":
				var webhook events.CreemRefundCreatedWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling refund.created webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}
				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.refund.created.webhook",
					OrderingKey: webhook.Object.Subscription.ID,
					Event:       webhook,
				}); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error publishing event",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				obs.Logger.InfoContext(spanCtx, "creemRefundCreatedWebhook event published", "webhook", webhook)
				return http.StatusOK, nil

			case "dispute.created":
				var webhook events.CreemDisputeCreatedWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling dispute.created webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}
				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.dispute.created.webhook",
					OrderingKey: webhook.Object.Subscription.ID,
					Event:       webhook,
				}); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error publishing event",
						Detail: err.Error(),
						Status: http.StatusInternalServerError,
					}
				}
				obs.Logger.InfoContext(spanCtx, "creemDisputeCreatedWebhook event published", "webhook", webhook)
				return http.StatusOK, nil

			default:
				obs.Logger.Warn("unsupported_creem_event_type",
					"eventType", payload.EventType,
					"eventID", payload.ID,
				)
				return nil, fuego.HTTPError{
					Title:  "unsupported event type",
					Detail: "event type " + payload.EventType + " is not supported",
					Status: http.StatusBadRequest,
				}
			}
		},
		option.Summary("creemWebhook"),
		option.Tags("webhooks"),
		option.Middleware(validateCreemWebhookSecretMiddleware))
}
