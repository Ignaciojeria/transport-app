package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderHeaders(c domain.Headers) table.OrderHeaders {
	return table.OrderHeaders{
		ID:          c.ID,
		Commerce:    c.Commerce,
		Consumer:    c.Consumer,
		ReferenceID: c.ReferenceID(),
		Organization: table.Organization{
			ID:      c.Organization.ID,
			Country: c.Organization.Country.Alpha2(),
		},
	}
}
