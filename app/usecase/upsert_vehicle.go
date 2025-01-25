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
		tidbrepository.NewUpsertVehicleHeaders,
		tidbrepository.NewUpsertVehicleCategory,
		tidbrepository.NewUpsertCarrier,
		tidbrepository.NewUpsertVehicle)
}

func NewUpsertVehicle(
	upsertVehicleHeaders tidbrepository.UpsertVehicleHeaders,
	upsertVehicleCategory tidbrepository.UpsertVehicleCategory,
	upsertCarrier tidbrepository.UpsertCarrier,
	upsertVehicle tidbrepository.UpsertVehicle,
) UpsertVehicle {
	return func(ctx context.Context, vehicle domain.Vehicle) error {
		vehicle.Headers.Organization = vehicle.Organization
		vehicleHeaders, err := upsertVehicleHeaders(ctx, vehicle.Headers)
		if err != nil {
			return err
		}

		vehicle.VehicleCategory.Organization = vehicle.Organization
		category, err := upsertVehicleCategory(ctx, vehicle.VehicleCategory)
		if err != nil {
			return err
		}

		vehicle.Carrier.Organization = vehicle.Organization
		carrier, err := upsertCarrier(ctx, vehicle.Carrier)
		if err != nil {
			return err
		}
		vehicle.Headers = vehicleHeaders
		vehicle.VehicleCategory = category
		vehicle.Carrier = carrier
		_, err = upsertVehicle(ctx, vehicle)
		return err
	}
}
