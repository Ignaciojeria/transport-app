package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
)

func init() {
	ioc.Register(checkSubscription)
}

func checkSubscription(
	s httpserver.Server,
	hasActiveSubscription supabaserepo.HasActiveSubscription,
	obs observability.Observability,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware) {
	fuego.Get(s.Manager, "/check-subscription",
		func(c fuego.ContextNoBody) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "checkSubscription")
			defer span.End()

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

			// Verificar si tiene suscripción activa
			hasSubscription, err := hasActiveSubscription(spanCtx, userID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_checking_subscription", "error", err, "user_id", userIDStr)
				return nil, err
			}

			// Devolver el resultado en JSON
			return map[string]bool{
				"has_active_subscription": hasSubscription,
			}, nil
		}, option.Summary("checkSubscription"), option.Middleware(jwtAuthMiddleware))
}
