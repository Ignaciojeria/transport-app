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
	ioc.Registry(optimizePlan, httpserver.New, usecase.NewOptimizePlan)
}
func optimizePlan(s httpserver.Server, optimize usecase.OptimizePlan) {
	fuego.Post(s.Manager, "/plans/optimize",
		func(c fuego.ContextWithBody[request.OptimizePlanRequest]) (request.SearchPlanByReferenceResponse, error) {
			req, err := c.Body()
			if err != nil {
				return request.SearchPlanByReferenceResponse{}, err
			}
			plan, err := optimize(c.Context(), req.Map())
			return request.MapSearchPlanByReferenceResponse(plan), err
		},
		option.Tags(tagPlanning),
		option.Summary("optimizePlan"))
}
