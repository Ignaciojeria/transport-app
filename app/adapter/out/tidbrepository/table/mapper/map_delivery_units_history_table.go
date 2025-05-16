package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapDeliveryUnitsHistoryTable(ctx context.Context, p domain.Plan) []table.DeliveryUnitsHistory {
	var deliveryUnitsHistory []table.DeliveryUnitsHistory

	for _, route := range p.Routes {
		for _, order := range route.Orders {
			for _, pkg := range order.Packages {
				deliveryUnitsHistory = append(deliveryUnitsHistory, table.DeliveryUnitsHistory{
					OrderDoc:             string(order.DocID(ctx)),
					DeliveryUnitDoc:      string(pkg.DocID(ctx, string(order.ReferenceID))),
					RouteDoc:             string(route.DocID(ctx)),
					VehicleDoc:           string(route.Vehicle.DocID(ctx)),
					CarrierDoc:           string(route.Vehicle.Carrier.DocID(ctx)),
					DriverDoc:            string(route.Vehicle.Carrier.Driver.DocID(ctx)),
					OrderStatusDoc:       pkg.Status.DocID().String(),
					NonDeliveryReasonDoc: string(pkg.ConfirmDelivery.NonDeliveryReason.DocID(ctx)),
				})
			}

		}
	}

	return nil
}
