package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapCarrierToTable(e domain.Carrier) table.Carrier {
	return table.Carrier{
		ID:             e.ID,
		OrganizationID: e.Organization.ID,
		ReferenceID:    e.ReferenceID,
		Name:           e.Name,
		NationalID:     e.NationalID,
	}
}
