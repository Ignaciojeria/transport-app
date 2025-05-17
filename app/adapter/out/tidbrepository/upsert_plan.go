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

type UpsertPlan func(ctx context.Context, p domain.Plan) error

func init() {
	ioc.Registry(NewUpsertPlan, database.NewConnectionFactory)
}

func NewUpsertPlan(conn database.ConnectionFactory) UpsertPlan {
	return func(ctx context.Context, p domain.Plan) error {
		var existing table.Plan
		err := conn.DB.WithContext(ctx).
			Table("plans").
			Where("document_id = ?", p.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newPlan := mapper.MapPlanToTable(ctx, p)
			return conn.Create(&newPlan).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(p)
		if !changed {
			return nil
		}

		updateData := mapper.MapPlanToTable(ctx, updated)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Save(&updateData).Error
	}
}
