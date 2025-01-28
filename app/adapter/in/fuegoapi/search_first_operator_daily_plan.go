package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(searchFirstOperatorDailyPlan, httpserver.New)
}
func searchFirstOperatorDailyPlan(s httpserver.Server) {
	fuego.Get(s.Manager, "/operator/{referenceID}/daily-plan",
		func(c fuego.ContextNoBody) (response.SearchFirstOperatorDailyPlanResponse, error) {

			return response.SearchFirstOperatorDailyPlanResponse{}, nil
		},
		option.Summary("searchFirstOperatorDailyPlan"),
		option.Tags(tagEndToEndOperator),
	)
}
