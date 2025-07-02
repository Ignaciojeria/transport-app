package workers

import (
	"context"
	"transport-app/app/adapter/out/fuegoapiclient"
	"transport-app/app/adapter/out/vroom"
	"transport-app/app/domain/optimization"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type FleetOptimizer func(
	ctx context.Context,
	input optimization.FleetOptimization) error

func init() {
	ioc.Registry(
		NewFleetOptimizer,
		vroom.NewOptimize,
		fuegoapiclient.NewPostUpsertRoute,
	)
}

func NewFleetOptimizer(
	optimize vroom.Optimize,
	postUpsertRoute fuegoapiclient.PostUpsertRoute) FleetOptimizer {
	return func(ctx context.Context, input optimization.FleetOptimization) error {

		plan, err := optimize(ctx, input)
		if err != nil {
			return err
		}

		unassignedRoute, ok := plan.GetUnassignedRoute()
		if ok {
			plan.Routes = append(plan.Routes, unassignedRoute)
		}

		for _, route := range plan.Routes {
			err := postUpsertRoute(ctx, route)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
