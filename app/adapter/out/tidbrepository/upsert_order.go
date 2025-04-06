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

		if errors.Is(err, gorm.ErrRecordNotFound) {
			DBOrderToCreate := mapper.MapOrderToTable(o)
			// Use the same Omit pattern for consistency
			return conn.
				Create(&DBOrderToCreate).Error
		}

		orderWithChanges, _ := order.Map().UpdateIfChanged(o)
		DBOrderToUpdate := mapper.MapOrderToTable(orderWithChanges)
		// Preservar todos los IDs de documento actuales
		DBOrderToUpdate.ID = order.ID
		DBOrderToUpdate.OrderHeadersDoc = order.OrderHeadersDoc
		DBOrderToUpdate.OriginNodeInfoDoc = order.OriginNodeInfoDoc
		DBOrderToUpdate.DestinationNodeInfoDoc = order.DestinationNodeInfoDoc
		DBOrderToUpdate.OrderStatusDoc = order.OrderStatusDoc
		DBOrderToUpdate.OrderTypeDoc = order.OrderTypeDoc
		DBOrderToUpdate.OriginContactDoc = order.OriginContactDoc
		DBOrderToUpdate.DestinationContactDoc = order.DestinationContactDoc
		DBOrderToUpdate.OriginAddressInfoDoc = order.OriginAddressInfoDoc
		DBOrderToUpdate.DestinationAddressInfoDoc = order.DestinationAddressInfoDoc
		DBOrderToUpdate.CreatedAt = order.CreatedAt
		DBOrderToUpdate.RouteDoc = order.RouteDoc

		// Actualizar IDs de documento si han cambiado
		if o.Headers.DocID().ShouldUpdate(order.OrderHeadersDoc) {
			DBOrderToUpdate.OrderHeadersDoc = o.Headers.DocID().String()
		}

		if o.Origin.DocID().ShouldUpdate(order.OriginNodeInfoDoc) {
			DBOrderToUpdate.OriginNodeInfoDoc = o.Origin.DocID().String()
		}

		if o.Destination.DocID().ShouldUpdate(order.DestinationNodeInfoDoc) {
			DBOrderToUpdate.DestinationNodeInfoDoc = o.Destination.DocID().String()
		}

		if o.OrderStatus.DocID().ShouldUpdate(order.OrderStatusDoc) {
			DBOrderToUpdate.OrderStatusDoc = o.OrderStatus.DocID().String()
		}

		if o.OrderType.DocID().ShouldUpdate(order.OrderTypeDoc) {
			DBOrderToUpdate.OrderTypeDoc = o.OrderType.DocID().String()
		}

		if o.Origin.Contact.DocID().ShouldUpdate(order.OriginContactDoc) {
			DBOrderToUpdate.OriginContactDoc = o.Origin.Contact.DocID().String()
		}

		if o.Destination.Contact.DocID().ShouldUpdate(order.DestinationContactDoc) {
			DBOrderToUpdate.DestinationContactDoc = o.Destination.Contact.DocID().String()
		}

		if o.Origin.AddressInfo.DocID().ShouldUpdate(order.OriginAddressInfoDoc) {
			DBOrderToUpdate.OriginAddressInfoDoc = o.Origin.AddressInfo.DocID().String()
		}

		if o.Destination.AddressInfo.DocID().ShouldUpdate(order.DestinationAddressInfoDoc) {
			DBOrderToUpdate.DestinationAddressInfoDoc = o.Destination.AddressInfo.DocID().String()
		}

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
