package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderReferences(
	ctx context.Context, order domain.Order) []table.OrderReferences {
	mapped := make([]table.OrderReferences, len(order.References))
	for i, ref := range order.References {
		mapped[i] = table.OrderReferences{
			DocumentID: ref.DocID(ctx, order.ReferenceID.String()).String(),
			OrderDoc:   order.DocID(ctx).String(),
			Type:       ref.Type,
			Value:      ref.Value,
		}
	}
	return mapped
}
