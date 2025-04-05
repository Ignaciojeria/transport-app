package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapCarrierToTable(e domain.Carrier) table.Carrier {
	return table.Carrier{
		OrganizationID: e.Organization.ID,
		Name:           e.Name,
		NationalID:     e.NationalID,
	}
}
