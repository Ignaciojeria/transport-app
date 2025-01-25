package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapVehicleCategory(domain domain.VehicleCategory) table.VehicleCategory {
	return table.VehicleCategory{
		ID:                    domain.ID,
		OrganizationCountryID: domain.Organization.OrganizationCountryID,
		Type:                  domain.Type,
		MaxPackagesQuantity:   domain.MaxPackagesQuantity,
	}
}
