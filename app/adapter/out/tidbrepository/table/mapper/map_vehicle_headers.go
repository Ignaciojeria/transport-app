package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapVehicleHeaders(c domain.Headers) table.VehicleHeaders {
	return table.VehicleHeaders{
		ID:       c.ID,
		Commerce: c.Commerce,
		Consumer: c.Consumer,
		OrganizationCountry: table.OrganizationCountry{
			ID: c.Organization.OrganizationCountryID,
		},
	}
}
