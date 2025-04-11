package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/shared/sharedcontext"
)

func MapOrganizationToTable(ctx context.Context, orgName string) table.Organization {
	return table.Organization{
		//ID:      org.ID,
		Name:    orgName,
		Country: sharedcontext.TenantCountryFromContext(ctx),
	}
}
