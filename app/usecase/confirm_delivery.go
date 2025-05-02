package usecase

import (
	"context"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type ConfirmDelivery func(ctx context.Context, input []domain.OrderHistory) error

func init() {
	ioc.Registry(NewConfirmDelivery)
}

func NewConfirmDelivery() ConfirmDelivery {
	return func(ctx context.Context, input []domain.OrderHistory) error {
		return nil
	}
}
