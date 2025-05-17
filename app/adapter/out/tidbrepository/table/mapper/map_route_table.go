package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapRouteTable(ctx context.Context, r domain.Route, planDoc string) table.Route {
	return table.Route{
		ReferenceID:       r.ReferenceID,
		DocumentID:        string(r.DocID(ctx)),
		TenantID:          sharedcontext.TenantIDFromContext(ctx),
		OriginNodeInfoDoc: string(r.Origin.DocID(ctx)),
		VehicleDoc:        string(r.Vehicle.DocID(ctx)),
		DriverDoc:         string(r.Vehicle.Carrier.Driver.DocID(ctx)),
		CarrierDoc:        string(r.Vehicle.Carrier.DocID(ctx)),
		PlanDoc:           planDoc,
	}
}
