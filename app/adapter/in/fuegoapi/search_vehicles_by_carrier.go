package fuegoapi

import (
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
		searchVehiclesByCarrier,
		httpserver.New,
		usecase.NewSearchVehiclesByCarrier)
}
func searchVehiclesByCarrier(
	s httpserver.Server,
	search usecase.SearchVehiclesByCarrier) {
	fuego.Get(s.Manager, "/carriers/{referenceID}/vehicles",
		func(c fuego.ContextNoBody) ([]response.SearchVehiclesByCarrierResponse, error) {
			searchFilters := domain.VehicleSearchFilters{
				CarrierReferenceID: c.PathParam("referenceID"),
				Organization: domain.Organization{
					Key:     c.Header("organization-key"),
					Country: countries.ByName(c.Header("country")),
				},
			}
			vehicles, err := search(c.Context(), searchFilters)
			if err != nil {
				return nil, err
			}
			return response.MapSearchVehiclesByCarrierResponse(vehicles), nil
		},
		option.Summary("searchVehiclesByCarrier"),
		option.Header("organization-key", "api organization key", param.Required()),
		option.Tags(tagFleets),
	)
}
