package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchOrders func(ctx context.Context, input domain.OrderSearchFilters) ([]domain.Order, error)

func init() {
	ioc.Registry(
		NewSearchOrders,
		tidbrepository.NewSearchOrders)
}

func NewSearchOrders(search tidbrepository.SearchOrders) SearchOrders {
	return func(ctx context.Context, input domain.OrderSearchFilters) ([]domain.Order, error) {
		return search(ctx, input)
	}
}
