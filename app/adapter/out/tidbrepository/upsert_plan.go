package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertPlan func(context.Context, domain.Plan) error

func init() {
	ioc.Registry(
		NewUpsertPlan,
		database.NewConnectionFactory,
		NewLoadOrderStatuses)
}

func NewUpsertPlan(conn database.ConnectionFactory, loadOrderStatuses LoadOrderStatuses) UpsertPlan {
	return func(ctx context.Context, p domain.Plan) error {
		return nil
	}
}
