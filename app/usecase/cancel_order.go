package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/projection/deliveryunits"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CancelOrder func(ctx context.Context, input domain.Route) error

func init() {
	ioc.Registry(
		NewCancelOrder,
		deliveryunits.NewProjection,
		tidbrepository.NewFindDeliveryUnitsProjectionResult,
		tidbrepository.NewUpsertDeliveryUnitsHistory)
}

func NewCancelOrder(
	projection deliveryunits.Projection,
	findDeliveryUnitsProjectionResult tidbrepository.FindDeliveryUnitsProjectionResult,
	upsertDeliveryUnitsHistory tidbrepository.UpsertDeliveryUnitsHistory,
) CancelOrder {
	return func(ctx context.Context, input domain.Route) error {
		referenceIds := []string{}
		for _, order := range input.Orders {
			referenceIds = append(referenceIds, order.ReferenceID.String())
		}
		last := 100
		results, _, err := findDeliveryUnitsProjectionResult(ctx, domain.DeliveryUnitsFilter{
			Order: &domain.OrderFilter{
				ReferenceIds: referenceIds,
			},
			RequestedFields: map[string]any{
				projection.DeliveryUnit().String():      true,
				projection.ReferenceID().String():       true,
				projection.DeliveryUnitLPN().String():   true,
				projection.DeliveryUnitItems().String(): true,
			},
			OnlyLatestStatus: true,
			Pagination: domain.Pagination{
				Last: &last,
			},
		})
		if err != nil {
			return err
		}

		nonDeliveryReason := input.Orders[0].DeliveryUnits[0].ConfirmDelivery.NonDeliveryReason
		manualChange := input.Orders[0].DeliveryUnits[0].ConfirmDelivery.ManualChange
		input.Orders[0].DeliveryUnits = nil
		// Primero asignamos los LPNs y SKUs a las órdenes
		for i := range input.Orders {
			for _, deliveryUnit := range results {
				if deliveryUnit.OrderReferenceID == input.Orders[i].ReferenceID.String() {
					input.Orders[i].DeliveryUnits = append(input.Orders[i].DeliveryUnits, domain.DeliveryUnit{
						Lpn:    deliveryUnit.LPN,
						Items:  deliveryUnit.JSONItems.Map(),
						Status: domain.Status{Status: domain.StatusCancelled},
						ConfirmDelivery: domain.ConfirmDelivery{
							NonDeliveryReason: nonDeliveryReason,
							ManualChange:      manualChange,
						},
					})
				}
			}
		}

		// Después de asignar los LPNs y SKUs, asignamos los índices y actualizamos el estado
		for i := range input.Orders {
			input.Orders[i].AssignIndexesIfNoLPN()
		}

		return upsertDeliveryUnitsHistory(ctx, domain.Plan{
			Routes: []domain.Route{input},
		})
	}
}
