package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/usecase/optimization"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertPlan func(context.Context, domain.Plan) (domain.Plan, error)

func init() {
	ioc.Registry(
		NewUpsertPlan,
		tidbrepository.NewUpsertPlanType,
		tidbrepository.NewUpsertPlanningStatus,
		tidbrepository.NewUpsertOperator,
		tidbrepository.NewUpsertPlan,
		optimization.NewOptimize,
		NewCreateOrder)
}

func NewUpsertPlan(
	upsertPlanType tidbrepository.UpsertPlanType,
	upsertPlanningStatus tidbrepository.UpsertPlanningStatus,
	upsertOperator tidbrepository.UpsertOperator,
	upsertPlan tidbrepository.UpsertPlan,
	optimize optimization.Optimize,
	upsertOrder CreateOrder,
) UpsertPlan {
	return func(ctx context.Context, plan domain.Plan) (domain.Plan, error) {
		/*
			plan.PlanType.Organization = plan.Organization
			err := upsertPlanType(ctx, plan.PlanType)

			for _, order := range plan.UnassignedOrders {
				order.Organization = plan.Organization
				err := upsertOrder(ctx, order)
				if err != nil {
					return domain.Plan{}, err
				}
			}

			for _, route := range plan.Routes {
				for _, order := range route.Orders {
					order.Organization = plan.Organization
					err := upsertOrder(ctx, order)
					if err != nil {
						return domain.Plan{}, err
					}
				}
			}

			if err != nil {
				return domain.Plan{}, err
			}
			plan.PlanningStatus.Organization = plan.Organization
			err = upsertPlanningStatus(ctx, plan.PlanningStatus)
			if err != nil {
				return domain.Plan{}, err
			}
			for index := range plan.Routes {
				plan.Routes[index].Organization = plan.Organization
				plan.Routes[index].Operator.Organization = plan.Organization
				operator, err := upsertOperator(ctx, plan.Routes[index].Operator)
				if err != nil {
					return domain.Plan{}, err
				}
				plan.Routes[index].Operator = operator
				plan.Routes[index].Organization = plan.Organization
			}*/
		return domain.Plan{}, upsertPlan(ctx, plan)
	}
}
