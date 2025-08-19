package mapper

import (
	"context"
	"encoding/json"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapRouteTable(ctx context.Context, r domain.Route, contract interface{}, planDoc string) table.Route {
	// Convert contract to json.RawMessage
	rawContract, _ := json.Marshal(contract)

	return table.Route{
		ReferenceID:       r.ReferenceID,
		Raw:               rawContract,
		DocumentID:        string(r.DocID(ctx)),
		TenantID:          sharedcontext.TenantIDFromContext(ctx),
		EndNodeInfoDoc:    string(r.Destination.DocID(ctx)),
		OriginNodeInfoDoc: string(r.Origin.DocID(ctx)),
		VehicleDoc:        string(r.Vehicle.DocID(ctx)),
		DriverDoc:         string(r.Vehicle.Carrier.Driver.DocID(ctx)),
		CarrierDoc:        string(r.Vehicle.Carrier.DocID(ctx)),
		PlanDoc:           planDoc,
	}
}
