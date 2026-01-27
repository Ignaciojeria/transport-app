package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"
	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
)

func init() {
	ioc.Registry(
		getCredits,
		httpserver.New,
		observability.NewObservability,
		supabaserepo.NewGetUserCredits,
		supabaserepo.NewGetCreditTransactions,
		apimiddleware.NewJWTAuthMiddleware,
	)
}

func getCredits(
	s httpserver.Server,
	obs observability.Observability,
	getUserCredits supabaserepo.GetUserCredits,
	getCreditTransactions supabaserepo.GetCreditTransactions,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware,
) {
	fuego.Get(s.Manager, "/credits",
		func(c fuego.ContextNoBody) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "getCredits")
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

			// Obtener créditos del usuario
			userCredits, err := getUserCredits(spanCtx, userID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_getting_user_credits", "error", err, "user_id", userIDStr)
				return nil, err
			}

			// Obtener últimas transacciones (opcional, últimas 10)
			transactions, err := getCreditTransactions(spanCtx, userID, 10)
			if err != nil {
				obs.Logger.WarnContext(spanCtx, "error_getting_credit_transactions", "error", err, "user_id", userIDStr)
				// No fallar si no se pueden obtener las transacciones
				transactions = []billing.CreditTransaction{}
			}

			// Devolver respuesta
			return map[string]interface{}{
				"balance":     userCredits.Balance,
				"transactions": transactions,
			}, nil
		}, option.Summary("getCredits"), option.Middleware(jwtAuthMiddleware))
}
