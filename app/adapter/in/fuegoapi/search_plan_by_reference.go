package fuegoapi

import (
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/shared/infrastructure/httpserver"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
)

func init() {
	ioc.Registry(searchPlanByReference, httpserver.New)
}
func searchPlanByReference(s httpserver.Server) {
	fuego.Get(s.Manager, "/plans/{referenceID}",
		func(c fuego.ContextNoBody) (response.SearchPlanByReferenceResponse, error) {
			return response.SearchPlanByReferenceResponse{}, nil
		},
		option.Summary("searchPlanByReference"),
		option.Tags(tagPlanning))
}
