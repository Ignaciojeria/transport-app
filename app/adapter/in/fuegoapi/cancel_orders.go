package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(cancelOrders, httpserver.New)
}
func cancelOrders(s httpserver.Server) {
	fuego.Post(s.Manager, "/orders/cancel",
		func(c fuego.ContextWithBody[request.CancelOrdersRequest]) (any, error) {

			return "unimplemented", nil
		},
		option.Summary("cancel orders"),
		option.Header("tenant", "api tenant", param.Required()),
		option.Header("channel", "api channel", param.Required()),
		option.Tags(tagOrders),
	)
}
