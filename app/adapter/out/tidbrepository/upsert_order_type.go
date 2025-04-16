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

type UpsertOrderType func(context.Context, domain.OrderType) error

func init() {
	ioc.Registry(NewUpsertOrderType, database.NewConnectionFactory)
}
func NewUpsertOrderType(conn database.ConnectionFactory) UpsertOrderType {
	return func(ctx context.Context, ot domain.OrderType) error {
		var existing table.OrderType

		err := conn.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", ot.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil {
			// Ya existe, no hacer nada
			return nil
		}

		// No existe - crear nuevo registro
		newRecord := mapper.MapOrderType(ctx, ot) // Usando el mapper correcto

		return conn.DB.WithContext(ctx).
			Omit("Organization").
			Create(&newRecord).Error
	}
}
