package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapProvinceTable(ctx context.Context, p domain.Province) table.Province {
	return table.Province{
		Name:       p.String(),
		DocumentID: p.DocID(ctx).String(),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
	}
}
