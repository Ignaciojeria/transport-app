package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(searchPlanByReference, httpserver.New, usecase.NewSearchPlan)
}
func searchPlanByReference(
	s httpserver.Server,
	searchPlan usecase.SearchPlan) {
	fuego.Get(s.Manager, "/plans/{referenceID}",
		func(c fuego.ContextNoBody) (request.SearchPlanByReferenceResponse, error) {
			searchFilters := domain.OrderSearchFilters{}
			searchFilters.Organization.Key = c.Header("organization-key")
			searchFilters.Organization.Country = countries.ByName(c.Header("country"))
			searchFilters.PlanReferenceID = c.PathParam("referenceID")
			plan, err := searchPlan(c.Context(), searchFilters)
			if err != nil {
				return request.SearchPlanByReferenceResponse{}, nil
			}
			return request.MapSearchPlanByReferenceResponse(plan), nil
		},
		option.Summary("searchPlanByReference"),
		option.Header("organization-key", "api organization key", param.Required()),
		option.Header("country", "api country", param.Required()),
		option.Tags(tagPlanning))
}
