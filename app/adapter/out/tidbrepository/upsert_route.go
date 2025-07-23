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

type UpsertRoute func(context.Context, domain.Route, string, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertRoute, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertRoute(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertRoute {
	return func(ctx context.Context, r domain.Route, planDoc string, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.Route

			err := tx.WithContext(ctx).
				Table("routes").
				Where("document_id = ?", r.DocID(ctx)).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// No existe → insert
				newRoute := mapper.MapRouteTable(ctx, r, planDoc)
				if err := tx.Create(&newRoute).Error; err != nil {
					return err
				}

				// Persistir FSMState si está presente
				if len(fsmState) > 0 {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			// Ya existe → update
			updateData := mapper.MapRouteTable(ctx, r, planDoc)
			updateData.ID = existing.ID // necesario para que GORM haga UPDATE
			updateData.CreatedAt = existing.CreatedAt

			if err := tx.Save(&updateData).Error; err != nil {
				return err
			}

			// Persistir FSMState si está presente
			if len(fsmState) > 0 {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})
	}
}
