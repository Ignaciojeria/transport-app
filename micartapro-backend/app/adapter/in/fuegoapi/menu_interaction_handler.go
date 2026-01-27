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
		supabaserepo.NewConsumeCredits,
	)
}
func menuInteractionHandler(
	s httpserver.Server,
	obs observability.Observability,
	publisherManager eventprocessing.PublisherManager,
	idempotencyKeyMiddleware apimiddleware.IdempotencyKeyMiddleware,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
	getUserCredits supabaserepo.GetUserCredits,
	consumeCredits supabaserepo.ConsumeCredits,
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

			// Validar y consumir créditos antes de procesar la interacción
			// Cada uso del agente consume 1 crédito
			const creditsPerInteraction = 1
			
			// Verificar créditos disponibles
			userCredits, err := getUserCredits(spanCtx, userID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_getting_user_credits", "error", err, "user_id", userIDStr)
				return nil, fuego.HTTPError{
					Title:  "error checking credits",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			if userCredits.Balance < creditsPerInteraction {
				obs.Logger.WarnContext(spanCtx, "insufficient_credits",
					"user_id", userIDStr,
					"balance", userCredits.Balance,
					"required", creditsPerInteraction)
				return nil, fuego.HTTPError{
					Title:  "insufficient credits",
					Detail: "You don't have enough credits to use the agent. Please purchase credits to continue.",
					Status: http.StatusPaymentRequired,
				}
			}

			// Consumir créditos
			versionID := uuid.New().String()
			description := "Uso del agente de menú"
			_, err = consumeCredits(spanCtx, billing.ConsumeCreditsRequest{
				UserID:      userID,
				Amount:      creditsPerInteraction,
				Source:      "agent.usage",
				SourceID:    &versionID,
				Description: &description,
			})

			if err != nil {
				if err == supabaserepo.ErrInsufficientCredits {
					obs.Logger.WarnContext(spanCtx, "insufficient_credits_on_consume",
						"user_id", userIDStr,
						"balance", userCredits.Balance,
						"required", creditsPerInteraction)
					return nil, fuego.HTTPError{
						Title:  "insufficient credits",
						Detail: "You don't have enough credits to use the agent. Please purchase credits to continue.",
						Status: http.StatusPaymentRequired,
					}
				}
				obs.Logger.ErrorContext(spanCtx, "error_consuming_credits", "error", err, "user_id", userIDStr)
				return nil, fuego.HTTPError{
					Title:  "error consuming credits",
					Detail: err.Error(),
					Status: http.StatusInternalServerError,
				}
			}

			obs.Logger.InfoContext(spanCtx, "credits_consumed_successfully",
				"user_id", userIDStr,
				"amount", creditsPerInteraction,
				"versionID", versionID)

			// Generar version_id si no viene en el contexto
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
