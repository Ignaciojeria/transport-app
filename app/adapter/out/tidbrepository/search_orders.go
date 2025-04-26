package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchOrders func(context.Context, domain.OrderSearchFilters) ([]domain.Order, error)

func init() {
	ioc.Registry(NewSearchOrders, database.NewConnectionFactory)
}
func NewSearchOrders(conn database.ConnectionFactory) SearchOrders {
	return func(ctx context.Context, osf domain.OrderSearchFilters) ([]domain.Order, error) {
		return []domain.Order{
			{
				ReferenceID: "hello world",
			},
		}, nil
	}
}
