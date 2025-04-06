package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertPlan func(context.Context, domain.Plan) error

func init() {
	ioc.Registry(
		NewUpsertPlan,
		tidb.NewTIDBConnection,
		NewLoadOrderStatuses)
}

func NewUpsertPlan(conn tidb.TIDBConnection, loadOrderStatuses LoadOrderStatuses) UpsertPlan {
	return func(ctx context.Context, p domain.Plan) error {
		return nil
	}
}
