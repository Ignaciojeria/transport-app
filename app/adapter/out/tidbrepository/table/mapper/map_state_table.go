package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapStateTable(ctx context.Context, p domain.State) table.State {
	return table.State{
		Name:        p.String(),
		DocumentID:  p.DocID(ctx).String(),
		CountryCode: sharedcontext.TenantCountryFromContext(ctx),
	}
}
