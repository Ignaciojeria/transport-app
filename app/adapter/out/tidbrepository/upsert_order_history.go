package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertOrderHistory func(context.Context, domain.OrderHistory) error

func init() {
	ioc.Registry(NewUpsertOrderHistory, database.NewConnectionFactory)
}

func NewUpsertOrderHistory(conn database.ConnectionFactory) UpsertOrderHistory {
	return func(ctx context.Context, oh domain.OrderHistory) error {
		return nil
	}
}
