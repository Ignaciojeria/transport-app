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

type UpsertVehicleCategory func(context.Context, domain.VehicleCategory) error

func init() {
	ioc.Registry(NewUpsertVehicleCategory, database.NewConnectionFactory)
}

func NewUpsertVehicleCategory(conn database.ConnectionFactory) UpsertVehicleCategory {
	return func(ctx context.Context, vc domain.VehicleCategory) error {
		var existing table.VehicleCategory

		err := conn.DB.WithContext(ctx).
			Table("vehicle_categories").
			Where("document_id = ?", vc.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newVehicleCategory := mapper.MapVehicleCategory(ctx, vc)
			return conn.Omit("Tenant").Create(&newVehicleCategory).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(vc)
		if !changed {
			return nil // No hay cambios, no hacemos nada
		}

		updateData := mapper.MapVehicleCategory(ctx, updated)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Tenant").Save(&updateData).Error
	}
}
