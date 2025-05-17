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

type UpsertOrderReferences func(context.Context, domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertOrderReferences,
		database.NewConnectionFactory)
}
func NewUpsertOrderReferences(conn database.ConnectionFactory) UpsertOrderReferences {
	return func(ctx context.Context, order domain.Order) error {
		orderDocID := order.DocID(ctx)
		orderReferences := mapper.MapOrderReferences(ctx, order)

		return conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("order_doc = ?", orderDocID).
				Delete(&table.OrderReferences{}).Error; err != nil {
				return err
			}
			if len(orderReferences) > 0 {
				if err := tx.Save(&orderReferences).Error; err != nil {
					return err
				}
			}
			if len(orderReferences) == 0 {
				if err := tx.Create(&table.OrderReferences{

					OrderDoc: orderDocID.String(),
				}).Error; err != nil {
					return err
				}
			}
			return nil
		})
	}
}
