package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertDailyPlan func(context.Context, domain.Plan) (domain.Plan, error)

func init() {
	ioc.Registry(
		NewUpsertDailyPlan,
		tidbrepository.NewUpsertPlanType)
}

func NewUpsertDailyPlan(upsertPlanType tidbrepository.UpsertPlanType) UpsertDailyPlan {
	return func(ctx context.Context, plan domain.Plan) (domain.Plan, error) {
		plan.PlanType.Organization = plan.Organization
		planType, err := upsertPlanType(ctx, plan.PlanType)
		if err != nil {
			return domain.Plan{}, err
		}
		plan.PlanType = planType
		return domain.Plan{}, nil
	}
}
