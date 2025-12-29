package fuegoapi

import (
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
	ioc.Registry(userMenusInsertedWebhook, httpserver.New,
		observability.NewObservability,
		eventprocessing.NewPublisherStrategy,
		apimiddleware.NewValidateSupabaseWebhookSecretMiddleware,
	)
}
func userMenusInsertedWebhook(
	s httpserver.Server,
	obs observability.Observability,
	publisherManager eventprocessing.PublisherManager,
	validateSupabaseWebhookSecretMiddleware apimiddleware.ValidateSupabaseWebhookSecretMiddleware) {
	fuego.Post(s.Manager, "/webhooks/user-menus-inserted",
		func(c fuego.ContextWithBody[events.UserMenusInsertedWebhook]) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "userMenusInsertedWebhook")
			defer span.End()
			requestBody, _ := c.Body()
			spanCtx = sharedcontext.WithIdempotencyKey(spanCtx, requestBody.Record.MenuID)
			spanCtx = sharedcontext.WithUserID(spanCtx, requestBody.Record.UserID)
			requestBody.CreatedAtToISO8601()
			obs.Logger.Info("userMenusInsertedWebhook request received", "requestBody", requestBody)

			if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
				Topic:       "micartapro.events",
				Source:      "micartapro.api.menu.interaction",
				OrderingKey: requestBody.Record.MenuID,
				Event:       requestBody,
			}); err != nil {
				return nil, fuego.HTTPError{
					Title:  "error publishing event",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx, "userMenusInsertedWebhook event published", "requestBody", requestBody)
			return http.StatusOK, nil
		},
		option.Summary("userMenusInsertedWebhook"),
		option.Tags("webhooks"),
		option.Middleware(validateSupabaseWebhookSecretMiddleware))
}
