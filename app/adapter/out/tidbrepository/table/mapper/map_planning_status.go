package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapPlanningStatus(ps domain.PlanningStatus) table.PlanningStatus {
	return table.PlanningStatus{
		ID:             ps.ID,
		OrganizationID: ps.Organization.ID,
		Name:           ps.Value,
	}
}
