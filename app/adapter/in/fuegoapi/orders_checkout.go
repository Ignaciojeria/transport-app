package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/gcppublisher"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(
		ordersCheckout,
		httpserver.New,
		gcppublisher.NewApplicationEvents)
}
func ordersCheckout(
	s httpserver.Server,
	outbox gcppublisher.ApplicationEvents) {
	fuego.Post(s.Manager, "/orders/checkout",
		func(c fuego.ContextWithBody[[]request.OrdersCheckoutRequest]) (any, error) {

			return "unimplemented", nil
		},
		option.Summary("ordersCheckout"),
		option.Tags(tagOrders),
		option.Tags(tagEndToEndOperator),
	)
}
