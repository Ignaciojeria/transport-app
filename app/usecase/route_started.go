package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/shared/projection/deliveryunits"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type RouteStarted func(ctx context.Context, input domain.Route) error

func init() {
	ioc.Registry(NewRouteStarted,
		deliveryunits.NewProjection,
		tidbrepository.NewFindDeliveryUnitsProjectionResult,
		tidbrepository.NewUpsertDeliveryUnitsHistory)
}

func NewRouteStarted(
	projection deliveryunits.Projection,
	findDeliveryUnitsProjectionResult tidbrepository.FindDeliveryUnitsProjectionResult,
	upsertDeliveryUnitsHistory tidbrepository.UpsertDeliveryUnitsHistory) RouteStarted {
	return func(ctx context.Context, input domain.Route) error {
		for i := range input.Orders {
			needsHydration := false
			for j := range input.Orders[i].DeliveryUnits {
				if input.Orders[i].DeliveryUnits[j].Lpn == "" && len(input.Orders[i].DeliveryUnits[j].Items) == 0 {
					needsHydration = true
					break
				}
			}

			if needsHydration {
				// Buscamos los datos de la orden
				last := 100
				results, _, err := findDeliveryUnitsProjectionResult(ctx, domain.DeliveryUnitsFilter{
					Order: &domain.OrderFilter{
						ReferenceIds: []string{input.Orders[i].ReferenceID.String()},
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
				// Asignamos los LPNs y SKUs a la orden
				input.Orders[i].DeliveryUnits = nil
				for _, deliveryUnit := range results {
					if deliveryUnit.OrderReferenceID == input.Orders[i].ReferenceID.String() {
						input.Orders[i].DeliveryUnits = append(input.Orders[i].DeliveryUnits, domain.DeliveryUnit{
							Status: domain.Status{Status: domain.StatusInTransit},
							Lpn:    deliveryUnit.LPN,
							Items:  deliveryUnit.JSONItems.Map(),
						})
					}
				}
			} else {
				// Si ya tiene LPN y SKUs, solo actualizamos el estado y ConfirmDelivery
				for j := range input.Orders[i].DeliveryUnits {
					input.Orders[i].DeliveryUnits[j].Status = domain.Status{Status: domain.StatusInTransit}
				}
			}

			input.Orders[i].AssignIndexesIfNoLPN()

		}
		return upsertDeliveryUnitsHistory(ctx, domain.Plan{
			Routes: []domain.Route{input},
		})
	}
}
