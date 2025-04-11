package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapPlanningStatus(ctx context.Context, ps domain.PlanningStatus) table.PlanningStatus {
	return table.PlanningStatus{
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),
		Name:           ps.Value,
	}
}
