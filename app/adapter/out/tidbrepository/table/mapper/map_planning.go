package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapPlan(plan domain.Plan) table.Plan {
	return table.Plan{
		ID:                   plan.ID,
		ReferenceID:          plan.ReferenceID,
		StartNodeReferenceID: string(plan.Origin.ReferenceID),
		JSONStartLocation: table.JSONPlanLocation{
			Latitude:  plan.Origin.AddressInfo.PlanLocation.Lat(),
			Longitude: plan.Origin.AddressInfo.PlanLocation.Lon(),
		},
		OrganizationCountryID: plan.Organization.OrganizationCountryID,
		PlannedDate:           plan.PlannedDate,
		PlanTypeID:            plan.PlanType.ID,
		PlanningStatusID:      plan.PlanningStatus.ID,
	}
}
