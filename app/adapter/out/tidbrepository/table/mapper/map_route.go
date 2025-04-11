package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapRoute(ctx context.Context, route domain.Route) table.Route {
	var accountID, vehicleID, carrierID *int64

	return table.Route{
		//ID:                 route.ID,
		EndNodeReferenceID: string(route.Destination.ReferenceID),
		JSONEndLocation: table.JSONPlanLocation{
			Longitude: route.Destination.AddressInfo.Location.Lon(),
			Latitude:  route.Destination.AddressInfo.Location.Lat(),
		},
		AccountID:      accountID,
		VehicleID:      vehicleID,
		CarrierID:      carrierID,
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),
	}
}
