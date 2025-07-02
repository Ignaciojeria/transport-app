package fuegoapi

import (
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
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
		func(c fuego.ContextWithBody[request.UpsertRouteRequest]) (any, error) {
			fmt.Println("upsert route controller call done! :D")
			return "unimplemented", nil
		},
		option.Summary("upsert route"),
		option.Tags(tagRoutes))
}
