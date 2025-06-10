package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(fixCoordinates, httpserver.New)
}
func fixCoordinates(s httpserver.Server) {
	fuego.Post(s.Manager, "/orders/destinations/fix",
		func(c fuego.ContextWithBody[request.OrderDestinationFixRequest]) (any, error) {
			return "unimplemented", nil
		},
		option.Summary("fix destination"),
		option.Tags(tagOrders))
}
