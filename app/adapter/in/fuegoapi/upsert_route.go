package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(upsertRoute, httpserver.New)
}
func upsertRoute(s httpserver.Server) {
	fuego.Post(s.Manager, "/routes",
		func(c fuego.ContextNoBody) (any, error) {

			return "unimplemented", nil
		},
		option.Summary("upsert route"),
		option.Tags(tagRoutes))
}
