package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapSizeCategory(ctx context.Context, sc domain.SizeCategory) table.SizeCategory {
	return table.SizeCategory{
		DocumentID: string(sc.DocumentID(ctx)),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
		Code:       sc.Code,
	}
}
