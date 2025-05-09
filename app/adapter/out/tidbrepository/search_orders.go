package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/projection/orders"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchOrders func(context.Context, domain.OrderFilterInput) (domain.OrderSearchResult, error)

func init() {
	ioc.Registry(
		NewSearchOrders,
		database.NewConnectionFactory,
		orders.NewProjection)
}
func NewSearchOrders(conn database.ConnectionFactory, projection orders.Projection) SearchOrders {
	return func(ctx context.Context, osf domain.OrderFilterInput) (domain.OrderSearchResult, error) {
		return domain.OrderSearchResult{
			Plan: domain.Plan{
				Routes: []domain.Route{},
			},
			HasNextPage: false,
			EndCursor:   new(string),
		}, nil
	}
}
