package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateOrder func(ctx context.Context, input domain.Order) (domain.Order, error)

func init() {
	ioc.Registry(
		NewCreateOrder,
		tidbrepository.NewSaveOrderQuery,
		tidbrepository.NewSaveOrder)
}

func NewCreateOrder(
	saveOrderQuery tidbrepository.SaveOrderQuery,
	saveOrder tidbrepository.SaveOrder,
) CreateOrder {
	return func(ctx context.Context, inOrder domain.Order) (domain.Order, error) {
		existingOrderDetails, err := saveOrderQuery(ctx, inOrder)
		if err != nil {
			return domain.Order{}, err
		}
		fmt.Println(existingOrderDetails)
		return saveOrder(ctx, existingOrderDetails, inOrder)
	}
}
