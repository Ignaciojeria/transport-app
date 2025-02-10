package optimization

import (
	"context"
	"errors"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type Optimize func(ctx context.Context, plan domain.Plan) (domain.Plan, error)

func init() {
	ioc.Registry(
		NewOptimize,
		newOneRouteOptimization)
}
func NewOptimize(
	oneRouteOptimization oneRouteOptimization,
) Optimize {
	return func(ctx context.Context, plan domain.Plan) (domain.Plan, error) {
		if len(plan.Routes) == 1 {
			return oneRouteOptimization(ctx, plan)
		}
		return domain.Plan{}, errors.New("unimplemented")
	}
}
