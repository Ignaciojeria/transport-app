package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		searchOrdersByUniqueReferences,
		httpserver.New,
		usecase.NewSearchOrders)
}
func searchOrdersByUniqueReferences(s httpserver.Server, search usecase.SearchOrders) {
	fuego.Post(s.Manager, "/order/search",
		func(c fuego.ContextWithBody[request.SearchOrdersByUniqueReferencesRequest]) ([]response.SearchOrdersResponse, error) {
			req, err := c.Body()
			if err != nil {
				return nil, err
			}
			searchFilters := req.Map()
			searchFilters.Organization.Key = c.Header("organization-key")
			searchFilters.Organization.Country = countries.ByName(c.Header("country"))
			orders, err := search(c.Context(), searchFilters)
			if err != nil {
				return nil, err
			}
			return response.MapSearchOrdersResponse(orders), nil
		},
		option.Summary("searchOrdersByUniqueReferences"),
		option.Header("organization-key", "api organization key", param.Required()),
		option.Header("country", "api country", param.Required()),
	)
}
