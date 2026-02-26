package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/eventprocessing"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/billing"
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
		supabaserepo.NewGetUserCredits,
	)
}
func menuInteractionHandler(
	s httpserver.Server,
	obs observability.Observability,
	publisherManager eventprocessing.PublisherManager,
	idempotencyKeyMiddleware apimiddleware.IdempotencyKeyMiddleware,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
	getUserCredits supabaserepo.GetUserCredits,
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

			// Extraer user_id del contexto
			userIDStr, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userIDStr == "" {
				return nil, fuego.UnauthorizedError{
					Title:  "user_id not found in context",
					Detail: "user_id not found in context",
					Status: 401,
				}
			}

			// Parsear userID
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return nil, fuego.BadRequestError{
					Title:  "invalid user_id",
					Detail: "invalid user_id format",
					Status: 400,
				}
			}

			// Verificar créditos mínimos antes de publicar (el consumo real se hace por operación en el subscriber)
			userCredits, err := getUserCredits(spanCtx, userID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_getting_user_credits", "error", err, "user_id", userIDStr)
				return nil, fuego.HTTPError{
					Title:  "error checking credits",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			if userCredits.Balance < billing.CreditsPerAgentUsage {
				obs.Logger.WarnContext(spanCtx, "insufficient_credits",
					"user_id", userIDStr,
					"balance", userCredits.Balance,
					"required", billing.CreditsPerAgentUsage)
				return nil, fuego.HTTPError{
					Title:  "insufficient credits",
					Detail: "You don't have enough credits to use the agent. Please purchase credits to continue.",
					Status: http.StatusPaymentRequired,
				}
			}

			// Generar version_id si no viene en el contexto
			versionID := uuid.New().String()
			if existingVersionID, ok := sharedcontext.VersionIDFromContext(spanCtx); ok && existingVersionID != "" {
				versionID = existingVersionID
			} else {
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
