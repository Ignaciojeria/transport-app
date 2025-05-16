package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapNodeType(ctx context.Context, nt domain.NodeType) table.NodeType {
	return table.NodeType{
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
		Value:      nt.Value,
		DocumentID: string(nt.DocID(ctx)),
	}
}
