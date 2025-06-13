package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(routeCreated, httpserver.New)
}
func routeCreated(s httpserver.Server) {
	fuego.Post(s.Manager, "/insert-your-custom-pattern-here",
		func(c fuego.ContextNoBody) (any, error) {

			return "unimplemented", nil
		}, option.Summary("routeCreated"))
}
