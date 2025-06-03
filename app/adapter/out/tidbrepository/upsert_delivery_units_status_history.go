package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertDeliveryUnitsHistory func(ctx context.Context, c domain.Plan) error

func init() {
	ioc.Registry(NewUpsertDeliveryUnitsHistory, database.NewConnectionFactory)
}

func NewUpsertDeliveryUnitsHistory(conn database.ConnectionFactory) UpsertDeliveryUnitsHistory {
	return func(ctx context.Context, c domain.Plan) error {
		// Get all delivery units history records for this plan
		deliveryUnitsHistory := mapper.MapDeliveryUnitsHistoryTable(ctx, c)

		// For each delivery unit history record
		for _, duh := range deliveryUnitsHistory {
			var existing table.DeliveryUnitsStatusHistory
			err := conn.DB.WithContext(ctx).
				Table("delivery_units_status_histories").
				Where("document_id = ?", duh.DocumentID).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			// No existe → insert
			if err := conn.Create(&duh).Error; err != nil {
				return err
			}
		}

		return nil
	}
}
