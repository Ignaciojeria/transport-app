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

type UpsertCarrier func(context.Context, domain.Carrier) error

func init() {
	ioc.Registry(NewUpsertCarrier, database.NewConnectionFactory)
}

func NewUpsertCarrier(conn database.ConnectionFactory) UpsertCarrier {
	return func(ctx context.Context, c domain.Carrier) error {
		var existing table.Carrier

		err := conn.DB.WithContext(ctx).
			Table("carriers").
			Where("document_id = ?", c.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newCarrier := mapper.MapCarrier(ctx, c)
			return conn.Omit("Tenant").Create(&newCarrier).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(c)
		if !changed {
			return nil // No hay cambios, no hacemos nada
		}

		updateData := mapper.MapCarrier(ctx, updated)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Tenant").Save(&updateData).Error
	}
}
