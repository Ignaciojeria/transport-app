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

type UpsertDriver func(context.Context, domain.Driver) error

func init() {
	ioc.Registry(NewUpsertDriver, database.NewConnectionFactory)
}

func NewUpsertDriver(conn database.ConnectionFactory) UpsertDriver {
	return func(ctx context.Context, d domain.Driver) error {
		var existing table.Driver

		err := conn.DB.WithContext(ctx).
			Table("drivers").
			Where("document_id = ?", d.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newDriver := mapper.MapDriver(ctx, d)
			return conn.Omit("Tenant").Create(&newDriver).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(d)
		if !changed {
			return nil // No hay cambios, no hacemos nada
		}

		updateData := mapper.MapDriver(ctx, updated)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Tenant").Save(&updateData).Error
	}
}
