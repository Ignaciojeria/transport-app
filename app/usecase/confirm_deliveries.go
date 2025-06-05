package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type ConfirmDeliveries func(ctx context.Context, input domain.Route) error

func init() {
	ioc.Registry(NewConfirmDeliveries, tidbrepository.NewUpsertDeliveryUnitsHistory)
}

func NewConfirmDeliveries(
	upsertDeliveryUnitsHistory tidbrepository.UpsertDeliveryUnitsHistory) ConfirmDeliveries {
	return func(ctx context.Context, input domain.Route) error {
		for i := range input.Orders {
			for j := range input.Orders[i].DeliveryUnits {
				input.Orders[i].AssignIndexesIfNoLPN()
				input.Orders[i].DeliveryUnits[j].UpdateStatusBasedOnNonDelivery()
			}
		}
		return upsertDeliveryUnitsHistory(ctx, domain.Plan{
			Routes: []domain.Route{input},
		})
	}
}
