package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	"github.com/biter777/countries"
)

func MapOrganizationFromTable(tableOrg table.Organization) domain.Organization {
	return domain.Organization{
		Name:    tableOrg.Name,
		Email:   tableOrg.Email,
		Country: countries.ByName(tableOrg.Country),
		Key:     "", // Agregar si tienes lógica para gestionar la clave
	}
}
