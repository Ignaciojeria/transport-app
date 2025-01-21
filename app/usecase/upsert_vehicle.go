package usecase

import (
	"context"
	"fmt"
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
		if err := upsert(ctx, v); err != nil {
			return err
		}
		fmt.Println("works")
		return nil
	}
}
