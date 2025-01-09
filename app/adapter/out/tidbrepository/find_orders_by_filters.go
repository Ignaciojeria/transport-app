package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewFindOrdersByFilters,
		tidb.NewTIDBConnection)
}

type FindOrdersByFilters func(
	context.Context,
	domain.OrderSearchFilters) ([]domain.Order, error)

func NewFindOrdersByFilters(
	conn tidb.TIDBConnection,
) FindOrdersByFilters {
	return func(ctx context.Context, o domain.OrderSearchFilters) ([]domain.Order, error) {
		return nil, nil
	}
}
