package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/restyclient"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/google/uuid"
)

func init() {
	ioc.Registry(
		creemCustomerPortal, httpserver.New,
		restyclient.NewGetCreemCustomerPortal,
		supabaserepo.NewGetCustomerID,
		observability.NewObservability,
		apimiddleware.NewJWTAuthMiddleware)
}

func creemCustomerPortal(
	s httpserver.Server,
	getCreemCustomerPortal restyclient.GetCreemCustomerPortal,
	getCustomerID supabaserepo.GetCustomerID,
	obs observability.Observability,
	jwtAuthMiddleware apimiddleware.JWTAuthMiddleware) {
	fuego.Get(s.Manager, "/customer-portal",
		func(c fuego.ContextNoBody) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "creemCustomerPortal")
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

			// Obtener customer_id desde Supabase
			customerID, err := getCustomerID(spanCtx, userID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_getting_customer_id", "error", err, "user_id", userIDStr)
				return nil, fuego.NotFoundError{
					Title:  "subscription not found",
					Detail: "no subscription found for this user",
					Status: 404,
				}
			}

			// Obtener el portal del consumidor desde Creem
			portalResponse, err := getCreemCustomerPortal(spanCtx, customerID)
			if err != nil {
				obs.Logger.ErrorContext(spanCtx, "error_getting_customer_portal", "error", err, "customer_id", customerID)
				return nil, err
			}

			// Devolver la URL del portal en JSON
			return map[string]string{
				"customer_portal_link": portalResponse.CustomerPortalLink,
			}, nil
		}, option.Summary("creemCustomerPortal"), option.Middleware(jwtAuthMiddleware))
}
