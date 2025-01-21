package usecase

import (
	"context"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertVehicle func(ctx context.Context, input interface{}) (interface{}, error)

func init() {
	ioc.Registry(NewUpsertVehicle)
}

func NewUpsertVehicle() UpsertVehicle {
	return func(ctx context.Context, input interface{}) (interface{}, error) {
		return input, nil
	}
}
