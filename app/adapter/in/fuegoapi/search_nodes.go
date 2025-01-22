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
		searchNodes,
		httpserver.New,
		usecase.NewSearchNodes)
}
func searchNodes(s httpserver.Server, search usecase.SearchNodes) {
	fuego.Get(s.Manager, "/node",
		func(c fuego.ContextWithBody[request.SearchNodesRequest]) ([]response.SearchNodesResponse, error) {
			req, err := c.Body()
			if err != nil {
				return nil, err
			}
			searchFilters := req.Map()
			searchFilters.Organization.Key = c.Header("organization-key")
			searchFilters.Organization.Country = countries.ByName(c.Header("country"))
			res, err := search(c.Context(), searchFilters)
			if err != nil {
				return nil, err
			}

			return response.MapSearchNodesResponse(res), nil
		},
		option.Summary("searchNodes"),
		option.Tags(tagNetwork),
		option.Header("organization-key", "api organization key", param.Required()),
		option.Header("country", "api country", param.Required()),
	)
}
