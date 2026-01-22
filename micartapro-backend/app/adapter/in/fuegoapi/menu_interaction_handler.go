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
	"github.com/go-fuego/fuego/param"
	"github.com/google/uuid"
)

func init() {
	ioc.Registry(
		menuInteractionHandler,
		httpserver.New,
		observability.NewObservability,
		eventprocessing.NewPublisherStrategy,
		apimiddleware.NewIdempotencyKeyMiddleware,
		apimiddleware.NewJWTAuthMiddleware,
	)
}
func menuInteractionHandler(
	s httpserver.Server,
	obs observability.Observability,
	publisherManager eventprocessing.PublisherManager,
	idempotencyKeyMiddleware apimiddleware.IdempotencyKeyMiddleware,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Post(s.Manager, "/menu/interaction",
		func(c fuego.ContextWithBody[events.MenuInteractionRequest]) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "menuInteractionHandler")
			defer span.End()

			body, err := c.Body()
			if err != nil {
				return nil, fuego.HTTPError{
					Title:  "error getting request body",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			if err := body.Validate(); err != nil {
				return nil, fuego.HTTPError{
					Title:  "error validating request body",
					Detail: err.Error(),
					Status: http.StatusBadRequest,
				}
			}

			// Generar version_id si no viene en el contexto
			var versionID string
			if existingVersionID, ok := sharedcontext.VersionIDFromContext(spanCtx); ok && existingVersionID != "" {
				versionID = existingVersionID
			} else {
				versionID = uuid.New().String()
				spanCtx = sharedcontext.WithVersionID(spanCtx, versionID)
			}

			if err := publisherManager.Publish(spanCtx, eventprocessing.PublishRequest{
				Topic:       "micartapro.events",
				Source:      "micartapro.api.menu.interaction",
				OrderingKey: body.MenuID,
				Event:       body,
			}); err != nil {
				return nil, fuego.HTTPError{
					Title:  "error publishing event",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}
			obs.Logger.InfoContext(spanCtx, "menuInteractionRequest published", "requestBody", body, "versionID", versionID)
			return map[string]string{"versionId": versionID}, nil
		},
		option.Summary("menuInteractionRequest"),
		option.Tags("agents"),
		option.Header("Idempotency-Key", "uuidv7", param.Required()),
		option.Middleware(idempotencyKeyMiddleware),
		option.Middleware(jwtAuthMiddleware))
}
