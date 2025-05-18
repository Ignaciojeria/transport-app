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

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newRecord := mapper.MapOrderType(ctx, ot)
			return conn.DB.WithContext(ctx).
				Omit("Organization").
				Create(&newRecord).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(ot)
		if !changed {
			return nil // No hay cambios, no hacemos nada
		}

		updateData := mapper.MapOrderType(ctx, updated)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.DB.WithContext(ctx).
			Omit("Organization").
			Save(&updateData).Error
	}
}
