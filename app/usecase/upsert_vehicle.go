package usecase

import (
	"context"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertVehicle func(ctx context.Context, vehicle domain.Vehicle) error

func init() {
	ioc.Registry(
		NewUpsertVehicle)
}

func NewUpsertVehicle() UpsertVehicle {
	return func(ctx context.Context, vehicle domain.Vehicle) error {
		return nil
	}
}
