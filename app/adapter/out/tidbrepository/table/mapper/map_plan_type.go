package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapPlanType(ctx context.Context, pt domain.PlanType) table.PlanType {
	return table.PlanType{
		Name:           pt.Value,
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),
	}
}
