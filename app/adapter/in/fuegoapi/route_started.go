package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(routeStarted, httpserver.New)
}
func routeStarted(s httpserver.Server) {
	fuego.Post(s.Manager, "/route/started",
		func(c fuego.ContextWithBody[request.RouteStartedRequest]) (any, error) {
			return "unimplemented", nil
		},
		option.Summary("route started"),
		option.Tags(tagRoutes))
}
