package fuegoapi

import (
	"micartapro/app/adapter/out/restyclient"
	"micartapro/app/shared/infrastructure/httpserver"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		creemCheckout, httpserver.New,
		restyclient.NewGetCreemCheckoutUrl,
		observability.NewObservability)
}
func creemCheckout(s httpserver.Server, getCreemCheckoutUrl restyclient.GetCreemCheckoutUrl, obs observability.Observability) {
	fuego.Get(s.Manager, "/checkout",
		func(c fuego.ContextNoBody) (any, error) {
			spanCtx, span := obs.Tracer.Start(c.Context(), "creemCheckout")
			defer span.End()
			checkoutResponse, err := getCreemCheckoutUrl(spanCtx)
			if err != nil {
				return nil, err
			}

			return c.Redirect(302, checkoutResponse.CheckoutURL)
		}, option.Summary("creemCheckout"))
}
