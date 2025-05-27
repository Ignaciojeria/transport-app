package mapper

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapDeliveryUnitsLabels(ctx context.Context, order domain.Order) []table.DeliveryUnitsLabels {
	var labels []table.DeliveryUnitsLabels
	for _, pkg := range order.DeliveryUnits {
		if len(pkg.Labels) == 0 {
			labels = append(labels, table.DeliveryUnitsLabels{
				DocumentID:      domain.Reference{}.DocID(ctx, pkg.DocID(ctx, order.ReferenceID.String()).String()).String(),
				OrderDoc:        order.DocID(ctx).String(),
				DeliveryUnitDoc: pkg.DocID(ctx, order.ReferenceID.String()).String(),
				Type:            "",
				Value:           "",
			})
			continue
		}
		for _, label := range pkg.Labels {
			labels = append(labels, table.DeliveryUnitsLabels{
				DocumentID:      label.DocID(ctx, pkg.DocID(ctx, order.ReferenceID.String()).String()).String(),
				Type:            label.Type,
				Value:           label.Value,
				OrderDoc:        order.DocID(ctx).String(),
				DeliveryUnitDoc: pkg.DocID(ctx, order.ReferenceID.String()).String(),
			})
		}
	}
	if len(order.DeliveryUnits) == 0 {
		emptyPkg := domain.DeliveryUnit{}
		deliveryUnitDoc := emptyPkg.DocID(ctx, order.ReferenceID.String()).String()
		labels = append(labels, table.DeliveryUnitsLabels{
			DocumentID:      domain.Reference{}.DocID(ctx, deliveryUnitDoc).String(),
			OrderDoc:        order.DocID(ctx).String(),
			DeliveryUnitDoc: deliveryUnitDoc,
			Type:            "",
			Value:           "",
		})
	}
	return labels
}
