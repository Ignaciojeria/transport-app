package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapVehicleCategory(ctx context.Context, domain domain.VehicleCategory) table.VehicleCategory {
	return table.VehicleCategory{
		OrganizationID:      sharedcontext.TenantIDFromContext(ctx),
		Type:                domain.Type,
		MaxPackagesQuantity: domain.MaxPackagesQuantity,
	}
}
