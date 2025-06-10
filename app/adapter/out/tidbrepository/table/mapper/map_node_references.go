package mapper

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapNodeReferences(
	ctx context.Context, node domain.NodeInfo) []table.NodeReferences {
	mapped := make([]table.NodeReferences, len(node.References))
	for i, ref := range node.References {
		docId := ref.DocID(ctx).String()
		fmt.Println("docId", docId)
		mapped[i] = table.NodeReferences{
			DocumentID: docId,
			NodeDoc:    node.DocID(ctx).String(),
			Type:       ref.Type,
			Value:      ref.Value,
		}
	}
	if len(mapped) == 0 {
		mapped = append(mapped, table.NodeReferences{
			DocumentID: domain.Reference{}.DocID(ctx).String(),
			NodeDoc:    node.DocID(ctx).String(),
			Type:       "", // o algo como "EMPTY"
			Value:      "",
		})
	}
	return mapped
}
