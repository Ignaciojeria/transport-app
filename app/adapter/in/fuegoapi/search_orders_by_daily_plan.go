package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
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
	fuego.Post(s.Manager, "/orders/daily-plan-search",
		func(c fuego.ContextWithBody[request.SearchOrdersByDailyPlanRequest]) ([]request.SearchOrdersResponse, error) {
			req, err := c.Body()
			if err != nil {
				return nil, err
			}
			searchFilters := req.Map()
			searchFilters.Organization.SetKey(c.Header("organization"))
			orders, err := search(c.Context(), searchFilters)
			if err != nil {
				return nil, err
			}
			return request.MapSearchOrdersResponse(orders), nil
		},
		option.Summary("searchOrdersByDailyPlan"),
		option.Tags(tagOrders, tagEndToEndOperator),
	)
}
