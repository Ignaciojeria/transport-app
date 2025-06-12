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

type UpsertPoliticalArea func(context.Context, domain.PoliticalArea) error

func init() {
	ioc.Registry(NewUpsertPoliticalArea, database.NewConnectionFactory)
}

func NewUpsertPoliticalArea(conn database.ConnectionFactory) UpsertPoliticalArea {
	return func(ctx context.Context, pa domain.PoliticalArea) error {
		var existing table.PoliticalArea

		err := conn.DB.WithContext(ctx).
			Table("political_areas").
			Where("document_id = ?", pa.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newRecord := mapper.MapPoliticalArea(ctx, pa)
			return conn.Omit("Tenant").Create(&newRecord).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(pa)
		if !changed {
			return nil // No hay cambios, no hacemos nada
		}

		updateData := mapper.MapPoliticalArea(ctx, updated)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Tenant").Save(&updateData).Error
	}
}
