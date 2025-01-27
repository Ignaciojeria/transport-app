package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(createPlan, httpserver.New)
}
func createPlan(s httpserver.Server) {
	fuego.Post(s.Manager, "/plan",
		func(c fuego.ContextWithBody[request.CreatePlanRequest]) (any, error) {

			return "unimplemented", nil
		},
		option.Summary("createPlan"),
		option.Tags(tagPlanning),
	)
}
