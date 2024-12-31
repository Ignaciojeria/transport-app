package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateTransportOrder func(ctx context.Context, input domain.TransportOrder) (domain.TransportOrder, error)

func init() {
	ioc.Registry(
		NewCreateTransportOrder,
		tidbrepository.NewSaveTransportOrder)
}

func NewCreateTransportOrder(saveTO tidbrepository.SaveTransportOrder) CreateTransportOrder {
	return func(ctx context.Context, to domain.TransportOrder) (domain.TransportOrder, error) {
		return saveTO(ctx, to)
	}
}
