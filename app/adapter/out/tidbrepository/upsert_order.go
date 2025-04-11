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
			Where("document_id = ?",
				o.DocID(ctx)).
			First(&order).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			DBOrderToCreate := mapper.MapOrderToTable(ctx, o)
			return conn.
				Create(&DBOrderToCreate).Error
		}

		orderWithChanges, _ := order.Map().UpdateIfChanged(o)
		DBOrderToUpdate := mapper.MapOrderToTable(ctx, orderWithChanges)
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
		if o.Headers.DocID(ctx).ShouldUpdate(order.OrderHeadersDoc) {
			DBOrderToUpdate.OrderHeadersDoc = o.Headers.DocID(ctx).String()
		}

		if o.OrderStatus.DocID().ShouldUpdate(order.OrderStatusDoc) {
			DBOrderToUpdate.OrderStatusDoc = o.OrderStatus.DocID().String()
		}

		if o.Origin.DocID(ctx).ShouldUpdate(order.OriginNodeInfoDoc) {
			DBOrderToUpdate.OriginNodeInfoDoc = o.Origin.DocID(ctx).String()
		}

		if o.Destination.DocID(ctx).ShouldUpdate(order.DestinationNodeInfoDoc) {
			DBOrderToUpdate.DestinationNodeInfoDoc = o.Destination.DocID(ctx).String()
		}

		if o.OrderStatus.DocID().ShouldUpdate(order.OrderStatusDoc) {
			DBOrderToUpdate.OrderStatusDoc = o.OrderStatus.DocID().String()
		}

		if o.OrderType.DocID(ctx).ShouldUpdate(order.OrderTypeDoc) {
			DBOrderToUpdate.OrderTypeDoc = o.OrderType.DocID(ctx).String()
		}

		originContactDoc := o.Origin.AddressInfo.Contact.DocID(ctx)
		if originContactDoc.ShouldUpdate(order.OriginContactDoc) {
			DBOrderToUpdate.OriginContactDoc = originContactDoc.String()
		}

		if o.Destination.AddressInfo.Contact.DocID(ctx).ShouldUpdate(order.DestinationContactDoc) {
			DBOrderToUpdate.DestinationContactDoc = o.Destination.AddressInfo.Contact.DocID(ctx).String()
		}

		if o.Origin.AddressInfo.DocID(ctx).ShouldUpdate(order.OriginAddressInfoDoc) {
			DBOrderToUpdate.OriginAddressInfoDoc = o.Origin.AddressInfo.DocID(ctx).String()
		}

		if o.Destination.AddressInfo.DocID(ctx).ShouldUpdate(order.DestinationAddressInfoDoc) {
			DBOrderToUpdate.DestinationAddressInfoDoc = o.Destination.AddressInfo.DocID(ctx).String()
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
			return tx.Omit("Organization").Save(&DBOrderToUpdate).Error
		}); err != nil {
			return err
		}
		return nil
	}
}
