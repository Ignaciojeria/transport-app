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
		tidbrepository.NewUpsertPlanType,
		tidbrepository.NewUpsertPlanningStatus,
		tidbrepository.NewFindOrdersByFilters,
		tidbrepository.NewUpsertOperator,
		tidbrepository.NewUpsertDailyPlan,
		tidbrepository.NewLoadOrganizationCountry)
}

func NewUpsertDailyPlan(
	upsertPlanType tidbrepository.UpsertPlanType,
	upsertPlanningStatus tidbrepository.UpsertPlanningStatus,
	findOrders tidbrepository.FindOrdersByFilters,
	upsertOperator tidbrepository.UpsertOperator,
	upsertDailyPlan tidbrepository.UpsertDailyPlan,
	loadOrganizationCountry tidbrepository.LoadOrganizationCountry,
) UpsertDailyPlan {
	return func(ctx context.Context, plan domain.Plan) (domain.Plan, error) {
		org, err := loadOrganizationCountry(ctx, plan.Organization)
		if err != nil {
			return domain.Plan{}, err
		}
		plan.Organization = org
		plan.PlanType.Organization = plan.Organization
		planType, err := upsertPlanType(ctx, plan.PlanType)
		if err != nil {
			return domain.Plan{}, err
		}
		plan.PlanningStatus.Organization = plan.Organization
		planningStatus, err := upsertPlanningStatus(ctx, plan.PlanningStatus)
		if err != nil {
			return domain.Plan{}, err
		}
		plan.Routes[0].Operator.Organization = plan.Organization
		operator, err := upsertOperator(ctx, plan.Routes[0].Operator)
		if err != nil {
			return domain.Plan{}, err
		}
		orders, err := findOrders(ctx, plan.GetOrderSearchFilters())
		if err != nil {
			return domain.Plan{}, err
		}
		plan.Routes[0].Orders = orders
		plan.Routes[0].Organization = plan.Organization
		plan.Routes[0].Operator = operator
		plan.Routes[0].ReferenceID = plan.ReferenceID
		plan.PlanType = planType
		plan.PlanningStatus = planningStatus
		return domain.Plan{}, upsertDailyPlan(ctx, plan)
	}
}
