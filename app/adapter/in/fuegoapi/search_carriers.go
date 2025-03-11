package fuegoapi

import (
	"strconv"
	"transport-app/app/adapter/in/fuegoapi/response"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/httpserver"
	"transport-app/app/usecase"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
)

func init() {
	ioc.Registry(
		searchCarriers,
		httpserver.New,
		usecase.NewSearchCarriers)
}
func searchCarriers(s httpserver.Server, search usecase.SearchCarriers) {
	fuego.Get(s.Manager, "/carriers",
		func(c fuego.ContextNoBody) ([]response.SearchCarriersResponse, error) {
			searchFilters := domain.CarrierSearchFilters{}
			searchFilters.Organization.SetKey(c.Header("organization"))

			page := c.QueryParam("page")
			if page == "" {
				page = "0"
			}
			p, err := strconv.Atoi(page)
			if err != nil {
				return nil, err
			}
			searchFilters.Page = p

			size := c.QueryParam("size")
			if size == "" {
				size = "10"
			}
			s, err := strconv.Atoi(size)
			if err != nil {
				return nil, err
			}
			searchFilters.Size = s

			carriers, err := search(c.Context(), searchFilters)
			if err != nil {
				return nil, err
			}
			return response.MapSearchCarriersResponse(carriers), nil
		},
		option.Summary("searchCarriers"),
		option.Header("organization", "api organization key", param.Required()),
		option.Query("page", "Page number", param.Default("0")),
		option.Query("size", "Page size", param.Default("10")),
		option.Tags(tagFleets),
	)
}
