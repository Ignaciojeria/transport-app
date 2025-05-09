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

type UpsertOrderDeliveryUnits func(context.Context, domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertOrderDeliveryUnits,
		database.NewConnectionFactory)
}
func NewUpsertOrderDeliveryUnits(conn database.ConnectionFactory) UpsertOrderDeliveryUnits {
	return func(ctx context.Context, order domain.Order) error {
		orderPackages := mapper.MapOrderDeliveryUnits(ctx, order)
		return conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("order_doc = ?", order.DocID(ctx)).
				Delete(&table.OrderDeliveryUnit{}).Error; err != nil {
				return err
			}
			if len(orderPackages) > 0 {
				if err := tx.Save(&orderPackages).Error; err != nil {
					return err
				}
			}

			if len(orderPackages) == 0 {
				pkg := domain.Package{}
				if err := tx.Save(&table.OrderDeliveryUnit{
					OrderDoc:        order.DocID(ctx).String(),
					DeliveryUnitDoc: pkg.DocID(ctx, order.ReferenceID.String()).String()}).Error; err != nil {
					return err
				}
			}

			return nil
		})
	}
}
