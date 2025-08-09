package usecase

import (
	"context"
	"transport-app/app/adapter/out/agents"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type VehiclesInputKeyNormalizationWorkflow func(ctx context.Context, input interface{}) (map[string]string, error)

func init() {
	ioc.Registry(
		NewVehiclesInputKeyNormalizationWorkflow,
		agents.NewVehicleFieldNamesNormalizer)
}

func NewVehiclesInputKeyNormalizationWorkflow(
	vehicleFieldNamesNormalizer agents.VehicleFieldNamesNormalizer) VehiclesInputKeyNormalizationWorkflow {
	return func(ctx context.Context, input interface{}) (map[string]string, error) {
		return vehicleFieldNamesNormalizer(ctx, input)
	}
}
