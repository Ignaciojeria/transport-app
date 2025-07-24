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

type UpsertVehicleHeaders func(context.Context, domain.Headers, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertVehicleHeaders, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertVehicleHeaders(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertVehicleHeaders {
	return func(ctx context.Context, h domain.Headers, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.VehicleHeaders

			err := tx.WithContext(ctx).
				Table("vehicle_headers").
				Where("document_id = ?", h.DocID(ctx)).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// No existe → insert
				newHeaders := mapper.MapVehicleHeaders(ctx, h)
				if err := tx.Omit("Tenant").Create(&newHeaders).Error; err != nil {
					return err
				}

				// Persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			// Ya existe → update solo si cambió algo
			updated, changed := existing.Map().UpdateIfChanged(h)
			if !changed {
				// No hay cambios, solo persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			updateData := mapper.MapVehicleHeaders(ctx, updated)
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
