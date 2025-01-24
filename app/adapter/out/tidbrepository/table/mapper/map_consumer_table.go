package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapConsumerToTable(c domain.Consumer) table.Consumer {
	return table.Consumer{
		ID:                    c.ID,
		Name:                  c.Value,
		OrganizationCountryID: c.Organization.OrganizationCountryID,
	}
}
