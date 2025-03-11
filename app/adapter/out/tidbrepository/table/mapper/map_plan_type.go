package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapPlanType(pt domain.PlanType) table.PlanType {
	return table.PlanType{
		ID:             pt.ID,
		Name:           pt.Value,
		OrganizationID: pt.Organization.ID,
	}
}
