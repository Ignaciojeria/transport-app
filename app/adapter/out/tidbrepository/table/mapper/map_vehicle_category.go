package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapVehicleCategory(ctx context.Context, vc domain.VehicleCategory) table.VehicleCategory {
	return table.VehicleCategory{
		DocumentID:          string(vc.DocID(ctx)),
		TenantID:            sharedcontext.TenantIDFromContext(ctx),
		Type:                vc.Type,
		MaxPackagesQuantity: vc.MaxPackagesQuantity,
	}
}
