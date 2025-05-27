package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertDeliveryUnitsLabels func(context.Context, domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertDeliveryUnitsLabels,
		database.NewConnectionFactory)
}

func NewUpsertDeliveryUnitsLabels(conn database.ConnectionFactory) UpsertDeliveryUnitsLabels {
	return func(ctx context.Context, order domain.Order) error {
		orderDocID := order.DocID(ctx)
		labels := mapper.MapDeliveryUnitsLabels(ctx, order)

		return conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("order_doc = ?", orderDocID).
				Delete(&table.DeliveryUnitsLabels{}).Error; err != nil {
				return err
			}
			if len(labels) > 0 {
				if err := tx.Save(&labels).Error; err != nil {
					return err
				}
			}
			if len(labels) == 0 {
				if err := tx.Create(&table.DeliveryUnitsLabels{
					DeliveryUnitDoc: orderDocID.String(),
				}).Error; err != nil {
					return err
				}
			}
			return nil
		})
	}
}
