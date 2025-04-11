package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapPlan(ctx context.Context, plan domain.Plan) table.Plan {
	return table.Plan{
		ReferenceID:          plan.ReferenceID,
		StartNodeReferenceID: string(plan.Origin.ReferenceID),
		JSONStartLocation: table.JSONPlanLocation{
			Latitude:  plan.Origin.AddressInfo.Location.Lat(),
			Longitude: plan.Origin.AddressInfo.Location.Lon(),
		},
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),
		PlannedDate:    plan.PlannedDate,
	}
}
