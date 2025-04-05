package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapRoute(route domain.Route) table.Route {
	var accountID, vehicleID, carrierID *int64

	// Verificar si Vehicle existe y tiene ID válido
	if route.Vehicle.ID != 0 {
		vehicleID = &route.Vehicle.ID
	}

	return table.Route{
		ID:                 route.ID,
		EndNodeReferenceID: string(route.Destination.ReferenceID),
		JSONEndLocation: table.JSONPlanLocation{
			Longitude: route.Destination.AddressInfo.Location.Lon(),
			Latitude:  route.Destination.AddressInfo.Location.Lat(),
		},
		ReferenceID:    route.ReferenceID,
		PlanID:         route.PlanID,
		AccountID:      accountID,
		VehicleID:      vehicleID,
		CarrierID:      carrierID,
		OrganizationID: route.Organization.ID,
	}
}
