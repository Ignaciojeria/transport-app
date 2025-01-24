package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapCommerceToTable(c domain.Commerce) table.Commerce {
	return table.Commerce{
		ID:                    c.ID,
		Name:                  c.Value,
		OrganizationCountryID: c.Organization.OrganizationCountryID,
	}
}
