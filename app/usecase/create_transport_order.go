package usecase

import (
	"context"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateTransportOrder func(ctx context.Context, input domain.TransportOrder) (domain.TransportOrder, error)

func init() {
	ioc.Registry(NewCreateTransportOrder)
}

func NewCreateTransportOrder() CreateTransportOrder {
	return func(ctx context.Context, input domain.TransportOrder) (domain.TransportOrder, error) {
		return input, nil
	}
}
