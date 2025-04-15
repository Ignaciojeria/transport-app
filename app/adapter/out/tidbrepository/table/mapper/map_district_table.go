package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapDistrictTable(ctx context.Context, p domain.District) table.District {
	return table.District{
		Name:        p.String(),
		DocumentID:  p.DocID(ctx).String(),
		CountryCode: sharedcontext.TenantCountryFromContext(ctx),
	}
}
