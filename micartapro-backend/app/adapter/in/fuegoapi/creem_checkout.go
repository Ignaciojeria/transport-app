package fuegoapi

import (
	"micartapro/app/adapter/in/fuegoapi/apimiddleware"
	"micartapro/app/adapter/out/restyclient"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Register(creemCheckout)
}
func creemCheckout(s httpserver.Server, getCreemCheckoutUrl restyclient.GetCreemCheckoutUrl, obs observability.Observability, jwtAuthMiddleware apimiddleware.JWTAuthMiddleware) {
	fuego.Get(s.Manager, "/checkout",
		func(c fuego.ContextNoBody) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "creemCheckout")
			defer span.End()

			// Extraer user_id del contexto
			userID, ok := sharedcontext.UserIDFromContext(spanCtx)
			if !ok || userID == "" {
				return nil, fuego.UnauthorizedError{
					Title:  "user_id not found in context",
					Detail: "user_id not found in context",
					Status: 401,
				}
			}

			checkoutResponse, err := getCreemCheckoutUrl(spanCtx, userID)
			if err != nil {
				return nil, err
			}

			// Devolver la URL en JSON para que el frontend pueda abrirla en una nueva pestaña
			return map[string]string{
				"checkout_url": checkoutResponse.CheckoutURL,
			}, nil
		}, option.Summary("creemCheckout"), option.Middleware(jwtAuthMiddleware))
}
