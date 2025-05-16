package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapOrderType(ctx context.Context, ot domain.OrderType) table.OrderType {
	return table.OrderType{
		Type:        ot.Type,
		Description: ot.Description,
		TenantID:    sharedcontext.TenantIDFromContext(ctx),
		DocumentID:  string(ot.DocID(ctx)),
	}
}
