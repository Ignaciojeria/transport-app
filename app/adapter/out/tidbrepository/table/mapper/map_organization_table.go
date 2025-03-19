package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrganizationToTable(org domain.Organization) table.Organization {
	return table.Organization{
		ID:      org.ID,
		Name:    org.Name,
		Country: org.Country.Alpha2(),
	}
}
