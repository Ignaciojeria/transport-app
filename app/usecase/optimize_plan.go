package usecase

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/usecase/optimization"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OptimizePlan func(ctx context.Context, input domain.Plan) (domain.Plan, error)

func init() {
	ioc.Registry(NewOptimizePlan, optimization.NewOptimize)
}

func NewOptimizePlan(optimize optimization.Optimize) OptimizePlan {
	return func(ctx context.Context, input domain.Plan) (domain.Plan, error) {
		return optimize(ctx, input)
	}
}
