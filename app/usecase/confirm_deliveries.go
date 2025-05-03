package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"golang.org/x/sync/errgroup"
)

type ConfirmDeliverires func(ctx context.Context, input []domain.OrderHistory) error

func init() {
	ioc.Registry(NewConfirmDeliveries, tidbrepository.NewUpsertOrderHistory)
}

func NewConfirmDeliveries(upsert tidbrepository.UpsertOrderHistory) ConfirmDeliverires {
	return func(ctx context.Context, input []domain.OrderHistory) error {
		group, groupCtx := errgroup.WithContext(ctx)

		for _, v := range input {
			// capturar variable para evitar cierre sobre el mismo valor en la goroutine
			order := v
			group.Go(func() error {
				return upsert(groupCtx, order)
			})
		}

		return group.Wait()
	}
}
