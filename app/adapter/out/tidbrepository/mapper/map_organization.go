package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrganizationToTable(org domain.Organization) table.Organization {
	return table.Organization{
		Name:    org.Name,
		Email:   org.Email,
		Country: org.Country.Alpha2(),
	}
}
