package usecase

import (
	"context"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OrdersCheckout func(ctx context.Context, input []domain.OrderCheckout) error

func init() {
	ioc.Registry(NewOrdersCheckout)
}

func NewOrdersCheckout() OrdersCheckout {
	return func(ctx context.Context, input []domain.OrderCheckout) error {
		return nil
	}
}
