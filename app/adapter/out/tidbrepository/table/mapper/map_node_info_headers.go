package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapNodeInfoHeaders(ctx context.Context, h domain.Headers) table.NodeInfoHeaders {
	return table.NodeInfoHeaders{
		Commerce:   h.Commerce,
		Consumer:   h.Consumer,
		Channel:    h.Channel,
		DocumentID: h.DocID(ctx).String(),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
	}
}
