package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/model"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(createTransportOrder, httpserver.New)
}
func createTransportOrder(s httpserver.Server) {
	fuego.Post(s.Manager, "/transport-order",
		func(c fuego.ContextWithBody[model.CreateTransportOrderRequest]) (model.CreateTransportOrderResponse, error) {
			return model.CreateTransportOrderResponse{}, nil
		}, option.Summary("createTransportOrder"))
}
