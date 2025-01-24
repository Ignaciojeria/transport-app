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

type UpsertOrderType func(context.Context, domain.OrderType) (domain.OrderType, error)

func init() {
	ioc.Registry(NewUpsertOrderType, tidb.NewTIDBConnection)
}

func NewUpsertOrderType(conn tidb.TIDBConnection) UpsertOrderType {
	return func(ctx context.Context, ot domain.OrderType) (domain.OrderType, error) {
		var orderType table.OrderType
		err := conn.DB.WithContext(ctx).
			Table("order_types").
			Where("type = ? AND organization_country_id = ?", ot.Type, ot.Organization.OrganizationCountryID).
			First(&orderType).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.OrderType{}, err
		}
		orderTypeWithChanges := orderType.Map().UpdateIfChanged(ot)
		DBOrderTypeToUpdate := mapper.MapOrderType(orderTypeWithChanges)
		DBOrderTypeToUpdate.CreatedAt = orderType.CreatedAt
		if err := conn.Save(&DBOrderTypeToUpdate).Error; err != nil {
			return domain.OrderType{}, err
		}
		return DBOrderTypeToUpdate.Map(), nil
	}
}
