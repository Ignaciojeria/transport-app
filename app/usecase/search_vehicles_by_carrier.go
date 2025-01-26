package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchVehiclesByCarrier func(
	ctx context.Context,
	input domain.VehicleSearchFilters) ([]domain.Vehicle, error)

func init() {
	ioc.Registry(NewSearchVehiclesByCarrier,
		tidbrepository.NewSearchVehiclesByCarrier)
}

func NewSearchVehiclesByCarrier(
	search tidbrepository.SearchVehiclesByCarrier) SearchVehiclesByCarrier {
	return func(ctx context.Context, input domain.VehicleSearchFilters) ([]domain.Vehicle, error) {
		return search(ctx, input)
	}
}
