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

type UpsertOrderHeaders func(context.Context, domain.Headers) error

func init() {
	ioc.Registry(NewUpsertOrderHeaders, tidb.NewTIDBConnection)
}
func NewUpsertOrderHeaders(conn tidb.TIDBConnection) UpsertOrderHeaders {
	return func(ctx context.Context, h domain.Headers) error {
		var orderHeaders table.OrderHeaders
		err := conn.DB.WithContext(ctx).
			Table("order_headers").
			Where("reference_id = ?", h.ReferenceID()).
			First(&orderHeaders).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err == nil {
			return nil
		}
		orderHeadersTbl := mapper.MapOrderHeaders(h)
		return conn.DB.WithContext(ctx).Create(&orderHeadersTbl).Error
	}
}
