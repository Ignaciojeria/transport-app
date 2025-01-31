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
)

func init() {
	ioc.Registry(
		searchOrdersByDailyPlan,
		httpserver.New,
		usecase.NewSearchOrders)
}
func searchOrdersByDailyPlan(
	s httpserver.Server,
	search usecase.SearchOrders) {
	fuego.Post(s.Manager, "/order/daily-plan-search",
		func(c fuego.ContextWithBody[request.SearchOrdersByDailyPlanRequest]) ([]response.SearchOrdersResponse, error) {
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
		option.Summary("searchOrdersByDailyPlan"),
		option.Tags(tagOrders, tagEndToEndOperator),
	)
}
