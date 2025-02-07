package usecase

import (
	"context"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type Checkout func(ctx context.Context, input []domain.Checkout) error

func init() {
	ioc.Registry(NewCheckout)
}

func NewCheckout() Checkout {
	return func(ctx context.Context, input []domain.Checkout) error {
		return nil
	}
}
