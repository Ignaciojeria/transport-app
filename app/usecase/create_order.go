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
		tidbrepository.NewSaveOrder)
}

func NewCreateOrder(saveTO tidbrepository.SaveOrder) CreateOrder {
	return func(ctx context.Context, to domain.Order) (domain.Order, error) {
		return saveTO(ctx, to)
	}
}
