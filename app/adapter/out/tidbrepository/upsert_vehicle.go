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

type UpsertVehicle func(context.Context, domain.Vehicle, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertVehicle, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertVehicle(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertVehicle {
	return func(ctx context.Context, v domain.Vehicle, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.Vehicle

			err := tx.WithContext(ctx).
				Table("vehicles").
				Where("document_id = ?", v.DocID(ctx)).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// No existe → insert
				newVehicle := mapper.MapVehicle(ctx, v)
				if err := tx.Omit("Tenant").Create(&newVehicle).Error; err != nil {
					return err
				}

				// Persistir FSMState si está presente
				if len(fsmState) > 0 {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			// Ya existe → update solo si cambió algo
			updated, changed := existing.Map().UpdateIfChanged(v)
			if !changed {
				// No hay cambios, solo persistir FSMState si está presente
				if len(fsmState) > 0 {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			updateData := mapper.MapVehicle(ctx, updated)
			updateData.ID = existing.ID // necesario para que GORM haga UPDATE
			updateData.CreatedAt = existing.CreatedAt

			if err := tx.Omit("Tenant").Save(&updateData).Error; err != nil {
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
