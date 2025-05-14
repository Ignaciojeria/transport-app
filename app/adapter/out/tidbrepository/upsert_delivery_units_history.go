package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertDeliveryUnitsHistory func(ctx context.Context, c domain.Plan) error

func init() {
	ioc.Registry(NewUpsertDeliveryUnitsHistory, database.NewConnectionFactory)
}

func NewUpsertDeliveryUnitsHistory(conn database.ConnectionFactory) UpsertDeliveryUnitsHistory {
	return func(ctx context.Context, c domain.Plan) error {
		return nil
	}
}
