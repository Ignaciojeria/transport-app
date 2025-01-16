package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateOrder func(ctx context.Context, input domain.Order) (domain.Order, error)

func init() {
	ioc.Registry(
		NewCreateOrder,
		tidbrepository.NewSaveOrderQuery,
		tidbrepository.NewSaveOrder,
		tidbrepository.NewLoadOrderStatuses)
}

func NewCreateOrder(
	saveOrderQuery tidbrepository.SaveOrderQuery,
	saveOrder tidbrepository.SaveOrder,
	loadOrderStatuses tidbrepository.LoadOrderStatuses,
) CreateOrder {
	return func(ctx context.Context, inOrder domain.Order) (domain.Order, error) {
		order, err := saveOrderQuery(ctx, inOrder)
		if err != nil {
			return domain.Order{}, err
		}
		order.OrderStatus = loadOrderStatuses().Available()
		order.HydrateOrder(inOrder)
		return saveOrder(ctx, order)
	}
}
