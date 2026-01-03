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

			// Identificar el tipo de evento y convertir al tipo correcto
			switch payload.EventType {
			case "subscription.trialing":
				var webhook events.CreemSubscriptionTrialingWebhook
				if err := json.Unmarshal(bodyBytes, &webhook); err != nil {
					return nil, fuego.HTTPError{
						Title:  "error unmarshaling subscription.trialing webhook",
						Detail: err.Error(),
						Status: http.StatusBadRequest,
					}
				}

				spanCtx = sharedcontext.WithIdempotencyKey(spanCtx, webhook.ID)
				userID := webhook.Object.Metadata.UserID
				if userID == "" {
					obs.Logger.Error("user_id_missing_in_webhook_metadata", "eventID", webhook.ID)
					return nil, fuego.HTTPError{
						Title:  "invalid webhook payload",
						Detail: "user_id missing in metadata",
						Status: http.StatusBadRequest,
					}
				}
				spanCtx = sharedcontext.WithUserID(spanCtx, userID)

				if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
					Topic:       "micartapro.events",
					Source:      "micartapro.api.creem.subscription.trialing.webhook",
					OrderingKey: webhook.Object.Customer.ID,
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
