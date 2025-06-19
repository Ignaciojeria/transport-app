package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(containerOptimization, httpserver.New)
}
func containerOptimization(s httpserver.Server) {
	fuego.Post(s.Manager, "/picking-and-dispatch",
		func(c fuego.ContextWithBody[request.ContainerOptimizationRequest]) (any, error) {

			return "unimplemented", nil
		}, option.Summary("picking & dispatch"), option.Tags("optimization"))
}
