package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/shared/sharedcontext"
)

func MapTenantTable(ctx context.Context, orgName string) table.Tenant {
	return table.Tenant{
		Name:    orgName,
		Country: sharedcontext.TenantCountryFromContext(ctx),
	}
}
