package mapper

import (
	"context"
	"fmt"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapDeliveryUnitsHistoryTable(ctx context.Context, p domain.Plan) []table.DeliveryUnitsStatusHistory {
	var deliveryUnitsHistory []table.DeliveryUnitsStatusHistory

	for _, route := range p.Routes {
		for _, order := range route.Orders {
			for _, pkg := range order.DeliveryUnits {
				deliveryUnitsHistory = append(deliveryUnitsHistory, table.DeliveryUnitsStatusHistory{
					OrderDoc:                     string(order.DocID(ctx)),
					TenantID:                     sharedcontext.TenantIDFromContext(ctx),
					DeliveryUnitDoc:              string(pkg.DocID(ctx)),
					RouteDoc:                     string(route.DocID(ctx)),
					Channel:                      sharedcontext.ChannelFromContext(ctx),
					VehicleDoc:                   string(route.Vehicle.DocID(ctx)),
					CarrierDoc:                   string(route.Vehicle.Carrier.DocID(ctx)),
					DriverDoc:                    string(route.Vehicle.Carrier.Driver.DocID(ctx)),
					PlanDoc:                      string(p.DocID(ctx)),
					DeliveryUnitStatusDoc:        string(pkg.Status.DocID()),
					NonDeliveryReasonReferenceID: pkg.ConfirmDelivery.NonDeliveryReason.ReferenceID,
					NonDeliveryReason:            pkg.ConfirmDelivery.NonDeliveryReason.Reason,
					NonDeliveryDetail:            pkg.ConfirmDelivery.NonDeliveryReason.Details,
					RecipientFullName:            pkg.ConfirmDelivery.Recipient.FullName,
					RecipientNationalID:          pkg.ConfirmDelivery.Recipient.NationalID,
					EvidencePhotos:               MapEvidencePhotosTable(ctx, pkg.ConfirmDelivery.EvidencePhotos),
					ConfirmDeliveryHandledAt:     pkg.ConfirmDelivery.HandledAt,
					ConfirmDeliveryLatitude:      pkg.ConfirmDelivery.Latitude,
					ConfirmDeliveryLongitude:     pkg.ConfirmDelivery.Longitude,
					ManualChangePerformedBy:      pkg.ConfirmDelivery.ManualChange.PerformedBy,
					ManualChangeReason:           pkg.ConfirmDelivery.ManualChange.Reason,
					DocumentID: domain.HashByTenant(
						ctx,
						string(order.DocID(ctx)),
						string(pkg.DocID(ctx)),
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
						pkg.ConfirmDelivery.ManualChange.PerformedBy,
						pkg.ConfirmDelivery.ManualChange.Reason,
					).String(),
				})
			}
		}
	}

	return deliveryUnitsHistory
}
