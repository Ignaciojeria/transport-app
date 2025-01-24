package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertOrder func(context.Context, domain.Order) (domain.Order, error)

func init() {
	ioc.Registry(NewUpsertOrder, tidb.NewTIDBConnection)
}
func NewUpsertOrder(conn tidb.TIDBConnection) UpsertOrder {
	return func(ctx context.Context, o domain.Order) (domain.Order, error) {
		tbl := mapper.MapOrderToTable(o)
		if err := conn.Transaction(func(tx *gorm.DB) error {
			return tx.
				Omit("OrganizationCountry").
				Omit("Commerce").
				Omit("Consumer").
				Omit("OrderStatus").
				Omit("OrderType").
				Omit("OriginContact").
				Omit("DestinationContact").
				Omit("OriginAddressInfo").
				Omit("DestinationAddressInfo").
				Omit("OriginNodeInfo").
				Omit("DestinationNodeInfo").
				Save(&tbl).Error
		}); err != nil {
			return domain.Order{}, err
		}
		return domain.Order{}, nil
	}
}
