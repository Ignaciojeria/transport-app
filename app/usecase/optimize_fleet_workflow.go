package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/vroom"
	"transport-app/app/domain/optimization"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OptimizeFleetWorkflow func(ctx context.Context, input optimization.FleetOptimization) ([]request.UpsertRouteRequest, error)

func init() {
	ioc.Registry(
		NewOptimizeFleetWorkflow,
		vroom.NewOptimize,
		observability.NewObservability,
	)
}

func NewOptimizeFleetWorkflow(
	optimize vroom.Optimize,
	obs observability.Observability,
) OptimizeFleetWorkflow {
	return func(ctx context.Context, input optimization.FleetOptimization) ([]request.UpsertRouteRequest, error) {
		routeRequests, err := optimize(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to optimize fleet: %w", err)
		}
		fmt.Printf("Optimización completada. Se generaron %d rutas.\n", len(routeRequests))
		return routeRequests, nil
	}
}
