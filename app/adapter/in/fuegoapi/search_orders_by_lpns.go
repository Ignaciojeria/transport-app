package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		searchOrdersByLpns,
		httpserver.New,
		usecase.NewSearchOrders)
}
func searchOrdersByLpns(s httpserver.Server, search usecase.SearchOrders) {
	fuego.Post(s.Manager, "/orders/lpns-search",
		func(c fuego.ContextWithBody[request.SearchOrdersByLpnsRequest]) ([]response.SearchOrdersResponse, error) {
			req, err := c.Body()
			if err != nil {
				return nil, err
			}
			searchFilters := req.Map()
			orders, err := search(c.Context(), searchFilters)
			if err != nil {
				return nil, err
			}
			return response.MapOrdersToSearchOrdersResponse(orders), nil
		},
		option.Summary("searchOrdersByLpns"),
		option.Header("organization", "api organization key", param.Required()),
		option.Tags(tagOrders),
	)
}
