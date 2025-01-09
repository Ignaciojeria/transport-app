package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/model"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(searchOrders, httpserver.New)
}
func searchOrders(s httpserver.Server) {
	fuego.Post(s.Manager, "/orders/search",
		func(c fuego.ContextWithBody[model.SearchOrdersRequest]) ([]model.SearchOrdersResponse, error) {
			return nil, nil
		}, option.Summary("searchOrders"))
}
