package fuegoapi

import (
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(planOptimization, httpserver.New)
}
func planOptimization(s httpserver.Server) {
	fuego.Post(s.Manager, "/plan/optimize",
		func(c fuego.ContextNoBody) (any, error) {

			return "unimplemented", nil
		}, option.Summary("planOptimization"), option.Tags(
			tagPlanning,
			tagEndToEndOperator))
}
