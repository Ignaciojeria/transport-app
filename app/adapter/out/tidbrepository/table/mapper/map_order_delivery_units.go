package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderDeliveryUnits(ctx context.Context, order domain.Order) []table.OrderDeliveryUnit {
	var orderPackages []table.OrderDeliveryUnit
	for _, p := range order.DeliveryUnits {
		orderPackages = append(orderPackages, table.OrderDeliveryUnit{
			OrderDoc:        order.DocID(ctx).String(),
			DeliveryUnitDoc: p.DocID(ctx, order.ReferenceID.String()).String(),
		})
	}
	return orderPackages
}
