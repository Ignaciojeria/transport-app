package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapOrderHeaders(ctx context.Context, c domain.Headers) table.OrderHeaders {
	return table.OrderHeaders{
		Commerce:       c.Commerce,
		Consumer:       c.Consumer,
		DocumentID:     string(c.DocID(ctx)),
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),
	}
}
