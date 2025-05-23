package mapper

import (
	"context"
	"fmt"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapDeliveryUnitsHistoryTable(ctx context.Context, p domain.Plan) []table.DeliveryUnitsHistory {
	var deliveryUnitsHistory []table.DeliveryUnitsHistory

	for _, route := range p.Routes {
		for _, order := range route.Orders {
			for _, pkg := range order.DeliveryUnits {
				deliveryUnitsHistory = append(deliveryUnitsHistory, table.DeliveryUnitsHistory{
					OrderDoc:                 string(order.DocID(ctx)),
					DeliveryUnitDoc:          string(pkg.DocID(ctx, string(order.ReferenceID))),
					RouteDoc:                 string(route.DocID(ctx)),
					Channel:                  sharedcontext.ChannelFromContext(ctx),
					VehicleDoc:               string(route.Vehicle.DocID(ctx)),
					CarrierDoc:               string(route.Vehicle.Carrier.DocID(ctx)),
					DriverDoc:                string(route.Vehicle.Carrier.Driver.DocID(ctx)),
					PlanDoc:                  string(p.DocID(ctx)),
					DeliveryUnitStatusDoc:    string(pkg.Status.DocID()),
					NonDeliveryReasonDoc:     string(pkg.ConfirmDelivery.NonDeliveryReason.DocID(ctx)),
					RecipientFullName:        pkg.ConfirmDelivery.Recipient.FullName,
					RecipientNationalID:      pkg.ConfirmDelivery.Recipient.NationalID,
					EvidencePhotos:           MapEvidencePhotosTable(ctx, pkg.ConfirmDelivery.EvidencePhotos),
					ConfirmDeliveryHandledAt: pkg.ConfirmDelivery.HandledAt,
					ConfirmDeliveryLatitude:  pkg.ConfirmDelivery.Latitude,
					ConfirmDeliveryLongitude: pkg.ConfirmDelivery.Longitude,
					DocumentID: domain.HashByTenant(
						ctx,
						string(order.DocID(ctx)),
						string(pkg.DocID(ctx, string(order.ReferenceID))),
						string(route.DocID(ctx)),
						string(p.DocID(ctx)),
						string(pkg.Status.DocID()),
						string(pkg.ConfirmDelivery.NonDeliveryReason.DocID(ctx)),
						pkg.ConfirmDelivery.Recipient.FullName,
						pkg.ConfirmDelivery.Recipient.NationalID,
						pkg.ConfirmDelivery.EvidencePhotos.DocID(ctx).String(),
						pkg.ConfirmDelivery.HandledAt.Format(time.RFC3339),
						fmt.Sprintf("%f", pkg.ConfirmDelivery.Latitude),
						fmt.Sprintf("%f", pkg.ConfirmDelivery.Longitude),
					).String(),
				})
			}

		}
	}

	return deliveryUnitsHistory
}
