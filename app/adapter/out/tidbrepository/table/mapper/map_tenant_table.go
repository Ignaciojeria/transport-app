package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapTenantTable(ctx context.Context, tenant domain.Tenant) table.Tenant {
	return table.Tenant{
		Name:    tenant.Name,
		Country: tenant.Country.Alpha2(),
	}
}
