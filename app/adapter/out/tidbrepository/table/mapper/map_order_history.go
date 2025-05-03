package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOderHistory(ctx context.Context, e domain.OrderHistory) table.OrdersHistory {
	return table.OrdersHistory{
		OrderDoc:      string(e.Order.DocID(ctx)),
		VehicleDoc:    string(e.Vehicle.DocID(ctx)),
		CarrierDoc:    string(e.Vehicle.Carrier.DocID(ctx)),
		OrderStatusID: e.OrderStatus.ID,
		RouteDoc:      string(e.Route.DocID(ctx)),
	}
}
