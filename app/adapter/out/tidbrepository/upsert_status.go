package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewUpsertStatus,
		database.NewConnectionFactory)
}

type UpsertStatus func(ctx context.Context, status domain.Status) error

func NewUpsertStatus(conn database.ConnectionFactory) UpsertStatus {
	return func(ctx context.Context, status domain.Status) error {
		record := table.Status{
			Status:     status.Status,
			DocumentID: status.DocID().String(),
		}

		return conn.DB.WithContext(ctx).Table("statuses").Create(&record).Error
	}
}
