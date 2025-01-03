package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrganizationFromTable(tableOrg table.Organization) domain.Organization {
	return domain.Organization{
		Name:  tableOrg.Name,
		Email: tableOrg.Email,
		Key:   "", // Agregar si tienes lógica para gestionar la clave
	}
}
