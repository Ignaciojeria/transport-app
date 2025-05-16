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

type UpsertVehicle func(context.Context, domain.Vehicle) error

func init() {
	ioc.Registry(NewUpsertVehicle, database.NewConnectionFactory)
}

func NewUpsertVehicle(conn database.ConnectionFactory) UpsertVehicle {
	return func(ctx context.Context, v domain.Vehicle) error {
		var existing table.Vehicle

		err := conn.DB.WithContext(ctx).
			Table("vehicles").
			Where("document_id = ?", v.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newVehicle := mapper.MapVehicle(ctx, v)
			return conn.Omit("Tenant").Create(&newVehicle).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(v)
		if !changed {
			return nil // No hay cambios, no hacemos nada
		}

		updateData := mapper.MapVehicle(ctx, updated)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Tenant").Save(&updateData).Error
	}
}
