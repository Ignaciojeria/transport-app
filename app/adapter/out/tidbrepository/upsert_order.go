package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertOrder func(context.Context, domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertOrder,
		tidb.NewTIDBConnection)
}
func NewUpsertOrder(conn tidb.TIDBConnection) UpsertOrder {
	return func(ctx context.Context, o domain.Order) error {
		var order table.Order
		err := conn.DB.WithContext(ctx).
			Table("orders").
			Preload("Organization").
			Where("document_id = ?",
				o.DocID()).
			First(&order).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		orderWithChanges, _ := order.Map().UpdateIfChanged(o)

		DBOrderToUpdate := mapper.MapOrderToTable(orderWithChanges)
		DBOrderToUpdate.CreatedAt = order.CreatedAt
		if err := conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(
				&table.OrderReferences{},
				"order_id = ?", DBOrderToUpdate.ID).Error; err != nil {
				return err
			}
			if err := tx.Delete(&table.OrderPackage{},
				"order_id = ?",
				DBOrderToUpdate.ID).Error; err != nil {
				return err
			}
			if err := tx.
				Omit("Organization").
				Omit("OrderHeaders").
				Omit("OrderStatus").
				Omit("OrderType").
				Omit("OriginContact").
				Omit("DestinationContact").
				Omit("OriginAddressInfo").
				Omit("DestinationAddressInfo").
				Omit("OriginNodeInfo").
				Omit("DestinationNodeInfo").
				Omit("Route").
				Save(&DBOrderToUpdate).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}
}
