package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapDistrictTable(ctx context.Context, d domain.District) table.District {
	return table.District{
		Name:       d.String(),
		DocumentID: d.DocID(ctx).String(),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
	}
}
