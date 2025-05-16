package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapDriver(ctx context.Context, d domain.Driver) table.Driver {
	return table.Driver{
		DocumentID: string(d.DocID(ctx)),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
		Name:       d.Name,
		NationalID: d.NationalID,
		Email:      d.Email,
	}
}
