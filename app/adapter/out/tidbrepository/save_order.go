package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewSaveOrder,
		tidb.NewTIDBConnection,
		NewLoadOrderStatuses)
}

type SaveOrder func(
	ctx context.Context, orderToCreate domain.Order) (domain.Order, error)

func NewSaveOrder(
	conn tidb.TIDBConnection,
	loadOrderSorderStatuses LoadOrderStatuses,
) SaveOrder {
	return func(ctx context.Context, orderToCreate domain.Order) (domain.Order, error) {
		orderTable := mapper.MapOrderToTable(orderToCreate)
		err := conn.Transaction(func(tx *gorm.DB) error {
			return tx.Save(&orderTable).Error
		})
		return domain.Order{}, err
	}
}
