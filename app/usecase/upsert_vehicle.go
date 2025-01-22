package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertVehicle func(ctx context.Context, vehicle domain.Vehicle) error

func init() {
	ioc.Registry(
		NewUpsertVehicle,
		tidbrepository.NewUpsertVehicleQuery,
		tidbrepository.NewUpsertVehicle)
}

func NewUpsertVehicle(
	query tidbrepository.UpsertVehicleQuery,
	upsert tidbrepository.UpsertVehicle) UpsertVehicle {
	return func(ctx context.Context, vehicle domain.Vehicle) error {
		v, err := query(ctx, vehicle)
		if err != nil {
			return err
		}
		v.UpdateIfChanged(vehicle)
		v.Organization = vehicle.Organization
		//v.BusinessIdentifiers = vehicle.BusinessIdentifiers
		if err := upsert(ctx, v); err != nil {
			return err
		}
		return nil
	}
}
