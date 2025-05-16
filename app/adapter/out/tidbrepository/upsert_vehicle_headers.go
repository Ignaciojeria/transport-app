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

type UpsertVehicleHeaders func(context.Context, domain.Headers) error

func init() {
	ioc.Registry(NewUpsertVehicleHeaders, database.NewConnectionFactory)
}

func NewUpsertVehicleHeaders(conn database.ConnectionFactory) UpsertVehicleHeaders {
	return func(ctx context.Context, h domain.Headers) error {
		var existing table.VehicleHeaders

		err := conn.DB.WithContext(ctx).
			Table("vehicle_headers").
			Where("document_id = ?", h.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newHeaders := mapper.MapVehicleHeaders(ctx, h)
			return conn.Omit("Tenant").Create(&newHeaders).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(h)
		if !changed {
			return nil // No hay cambios, no hacemos nada
		}

		updateData := mapper.MapVehicleHeaders(ctx, updated)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Tenant").Save(&updateData).Error
	}
}
