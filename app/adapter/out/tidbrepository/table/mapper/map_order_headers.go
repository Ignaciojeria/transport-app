package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapOrderHeaders(ctx context.Context, c domain.Headers) table.OrderHeaders {
	return table.OrderHeaders{
		Commerce:   c.Commerce,
		Consumer:   c.Consumer,
		Channel:    c.Channel,
		DocumentID: string(c.DocID(ctx)),
		TenantID:   sharedcontext.TenantIDFromContext(ctx),
	}
}
