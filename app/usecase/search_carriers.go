package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchCarriers func(ctx context.Context, filters domain.CarrierSearchFilters) ([]domain.Carrier, error)

func init() {
	ioc.Registry(
		NewSearchCarriers,
		tidbrepository.NewSearchCarriers)
}

func NewSearchCarriers(search tidbrepository.SearchCarriers) SearchCarriers {
	return func(ctx context.Context, filters domain.CarrierSearchFilters) ([]domain.Carrier, error) {
		return search(ctx, filters)
	}
}
