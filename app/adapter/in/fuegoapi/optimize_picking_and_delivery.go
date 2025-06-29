package fuegoapi
/*
import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(optimizePickingAndDelivery, httpserver.New)
}
func optimizePickingAndDelivery(s httpserver.Server) {
	fuego.Post(s.Manager, "/optimize/picking-and-delivery",
		func(c fuego.ContextWithBody[request.OptimizePickingAndDeliveryRequest]) (any, error) {

			return "unimplemented", nil
		}, option.Summary("optimize picking & delivery"), option.Tags("optimization"))
}
*/