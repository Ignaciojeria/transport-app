package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapVehicleHeaders(ctx context.Context, h domain.Headers) table.VehicleHeaders {
	return table.VehicleHeaders{
		DocumentID: string(h.DocID(ctx)),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
		Commerce:   h.Commerce,
		Consumer:   h.Consumer,
	}
}
