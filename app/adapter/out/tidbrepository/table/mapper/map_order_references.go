package mapper

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderReferences(
	ctx context.Context, order domain.Order) []table.OrderReferences {
	mapped := make([]table.OrderReferences, len(order.References))
	for i, ref := range order.References {
		docId := ref.DocID(ctx, order.DocID(ctx).String()).String()
		fmt.Println("docId", docId)
		mapped[i] = table.OrderReferences{
			DocumentID: docId,
			OrderDoc:   order.DocID(ctx).String(),
			Type:       ref.Type,
			Value:      ref.Value,
		}
	}
	if len(mapped) == 0 {
		mapped = append(mapped, table.OrderReferences{
			DocumentID: domain.Reference{}.DocID(ctx, order.DocID(ctx).String()).String(),
			OrderDoc:   order.DocID(ctx).String(),
			Type:       "", // o algo como "EMPTY"
			Value:      "",
		})
	}
	return mapped
}
