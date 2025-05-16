package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapCarrier(ctx context.Context, c domain.Carrier) table.Carrier {
	return table.Carrier{
		DocumentID: string(c.DocID(ctx)),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
		Name:       c.Name,
		NationalID: c.NationalID,
	}
}
