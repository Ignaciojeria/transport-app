package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapVehicleCategory(domain domain.VehicleCategory) table.VehicleCategory {
	return table.VehicleCategory{
		OrganizationID:      domain.Organization.ID,
		Type:                domain.Type,
		MaxPackagesQuantity: domain.MaxPackagesQuantity,
	}
}
