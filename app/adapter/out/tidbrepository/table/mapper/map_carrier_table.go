package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapCarrierToTable(ctx context.Context, e domain.Carrier) table.Carrier {
	return table.Carrier{
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),
		Name:           e.Name,
		NationalID:     e.NationalID,
	}
}
