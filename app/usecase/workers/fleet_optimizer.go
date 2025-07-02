package workers

import (
	"context"
	"transport-app/app/adapter/out/vroom"
	"transport-app/app/domain/optimization"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type FleetOptimizer func(ctx context.Context, input optimization.FleetOptimization) error

func init() {
	ioc.Registry(NewFleetOptimizer, vroom.NewOptimize)
}

func NewFleetOptimizer(optimize vroom.Optimize) FleetOptimizer {
	return func(ctx context.Context, input optimization.FleetOptimization) error {

		plan, err := optimize(ctx, input)

		unassignedRoute, ok := plan.GetUnassignedRoute()
		if ok {
			plan.Routes = append(plan.Routes, unassignedRoute)
		}

		if err != nil {
			return err
		}

		return nil
	}
}
