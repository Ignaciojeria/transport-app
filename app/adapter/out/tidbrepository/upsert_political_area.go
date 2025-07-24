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

type UpsertPoliticalArea func(context.Context, domain.PoliticalArea, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertPoliticalArea, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertPoliticalArea(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertPoliticalArea {
	return func(ctx context.Context, pa domain.PoliticalArea, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.PoliticalArea

			err := tx.WithContext(ctx).
				Table("political_areas").
				Where("document_id = ?", pa.DocID(ctx)).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// No existe → insert
				newRecord := mapper.MapPoliticalArea(ctx, pa)
				if err := tx.Omit("Tenant").Create(&newRecord).Error; err != nil {
					return err
				}

				// Persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			// Ya existe → update solo si cambió algo
			updated, changed := existing.Map().UpdateIfChanged(pa)
			if !changed {
				// No hay cambios, solo persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			updateData := mapper.MapPoliticalArea(ctx, updated)
			updateData.ID = existing.ID // necesario para que GORM haga UPDATE
			updateData.CreatedAt = existing.CreatedAt

			if err := tx.Omit("Tenant").Save(&updateData).Error; err != nil {
				return err
			}

			// Persistir FSMState si está presente
			if len(fsmState) > 0 && saveFSMTransition != nil {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})
	}
}
