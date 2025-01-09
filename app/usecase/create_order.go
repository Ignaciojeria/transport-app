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
		tidbrepository.NewEnsureOrganizationForCountry,
		tidbrepository.NewSaveOrder)
}

func NewCreateOrder(
	ensureOrganizationForCountry tidbrepository.EnsureOrganizationForCountry,
	saveOrder tidbrepository.SaveOrder) CreateOrder {
	return func(ctx context.Context, to domain.Order) (domain.Order, error) {
		if err := to.ValidateCollectAvailabilityDate(); err != nil {
			return domain.Order{}, err
		}
		if err := to.ValidatePromisedDate(); err != nil {
			return domain.Order{}, err
		}
		if err := ensureOrganizationForCountry(ctx, to.Organization); err != nil {
			return domain.Order{}, err
		}
		return saveOrder(ctx, to)
	}
}
