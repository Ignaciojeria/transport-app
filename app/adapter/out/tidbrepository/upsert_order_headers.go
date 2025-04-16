package tidbrepository

import (
	"context"
	"errors"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertOrderHeaders func(context.Context, domain.Headers) error

func init() {
	ioc.Registry(NewUpsertOrderHeaders, database.NewConnectionFactory)
}

func NewUpsertOrderHeaders(conn database.ConnectionFactory) UpsertOrderHeaders {
	return func(ctx context.Context, h domain.Headers) error {
		var orderHeaders table.OrderHeaders

		err := conn.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", h.DocID(ctx)).
			First(&orderHeaders).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			// Ya existe → no hacer nada
			return nil
		}

		orderHeadersTbl := mapper.MapOrderHeaders(ctx, h)

		return conn.DB.
			WithContext(ctx).
			Omit("Organization").
			Create(&orderHeadersTbl).Error
	}
}
