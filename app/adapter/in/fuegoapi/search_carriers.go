package fuegoapi

import (
	"strconv"
	"transport-app/app/adapter/in/fuegoapi/response"
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
	ioc.Registry(
		searchCarriers,
		httpserver.New,
		usecase.NewSearchCarriers)
}
func searchCarriers(s httpserver.Server, search usecase.SearchCarriers) {
	fuego.Get(s.Manager, "/carrier",
		func(c fuego.ContextNoBody) ([]response.SearchCarriersResponse, error) {
			searchFilters := domain.CarrierSearchFilters{
				Organization: domain.Organization{
					Key:     c.Header("organization-key"),
					Country: countries.ByName(c.Header("country")),
				},
			}

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
		option.Header("organization-key", "api organization key", param.Required()),
		option.Header("country", "api country", param.Required()),
		option.Query("page", "Page number", param.Default("0")),
		option.Query("size", "Page size", param.Default("10")),
		option.Tags(tagFleets),
	)
}
