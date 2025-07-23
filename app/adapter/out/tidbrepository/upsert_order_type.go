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

type UpsertOrderType func(context.Context, domain.OrderType, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertOrderType, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertOrderType(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertOrderType {
	return func(ctx context.Context, ot domain.OrderType, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.OrderType

			err := tx.WithContext(ctx).
				Table("order_types").
				Where("document_id = ?", ot.DocID(ctx)).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// No existe → insert
				newRecord := mapper.MapOrderType(ctx, ot)
				if err := tx.Omit("Organization").Create(&newRecord).Error; err != nil {
					return err
				}

				// Persistir FSMState si está presente
				if len(fsmState) > 0 {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			// Ya existe → update solo si cambió algo
			updated, changed := existing.Map().UpdateIfChanged(ot)
			if !changed {
				// No hay cambios, solo persistir FSMState si está presente
				if len(fsmState) > 0 {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			updateData := mapper.MapOrderType(ctx, updated)
			updateData.ID = existing.ID // necesario para que GORM haga UPDATE
			updateData.CreatedAt = existing.CreatedAt

			if err := tx.Omit("Organization").Save(&updateData).Error; err != nil {
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
