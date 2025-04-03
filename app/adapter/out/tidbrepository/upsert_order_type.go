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

type UpsertOrderType func(context.Context, domain.OrderType) error

func init() {
	ioc.Registry(NewUpsertOrderType, tidb.NewTIDBConnection)
}

func NewUpsertOrderType(conn tidb.TIDBConnection) UpsertOrderType {
	return func(ctx context.Context, ot domain.OrderType) error {
		var existing table.OrderType

		err := conn.DB.WithContext(ctx).
			Table("order_types").
			Preload("Organization").
			Where("document_id = ?", ot.DocID()).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insertar nuevo
			newRecord := mapper.MapOrderType(ot)
			return conn.DB.WithContext(ctx).
				Omit("Organization").
				Create(&newRecord).Error
		}

		// Existe → ver si cambió
		updated, changed := existing.Map().UpdateIfChanged(ot)
		if !changed {
			return nil
		}

		toUpdate := mapper.MapOrderType(updated)
		toUpdate.ID = existing.ID
		toUpdate.CreatedAt = existing.CreatedAt

		return conn.DB.WithContext(ctx).
			Omit("Organization").
			Save(&toUpdate).Error
	}
}
