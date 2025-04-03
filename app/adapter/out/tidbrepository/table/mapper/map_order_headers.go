package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderHeaders(c domain.Headers) table.OrderHeaders {
	return table.OrderHeaders{
		ID:             c.ID,
		Commerce:       c.Commerce,
		Consumer:       c.Consumer,
		DocumentID:     string(c.DocID()),
		OrganizationID: c.Organization.ID,
	}
}
