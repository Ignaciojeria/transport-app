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
		tidbrepository.NewFindOrdersByReferenceAndFilters)
}

func NewSearchOrders(search tidbrepository.FindOrdersByReferenceAndFilters) SearchOrders {
	return func(ctx context.Context, input domain.OrderSearchFilters) ([]domain.Order, error) {
		return search(ctx, input)
	}
}
