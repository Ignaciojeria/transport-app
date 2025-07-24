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

type UpsertContact func(ctx context.Context, c domain.Contact, fsmState ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertContact, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertContact(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertContact {
	return func(ctx context.Context, c domain.Contact, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.Contact
			err := tx.WithContext(ctx).
				Table("contacts").
				Where("document_id = ?", c.DocID(ctx)).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// No existe → insert
				newContact := mapper.MapContactToTable(ctx, c)
				if err := tx.Omit("Tenant").Create(&newContact).Error; err != nil {
					return err
				}

				// Persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			// Ya existe → update solo si cambió algo
			updated, changed := existing.Map().UpdateIfChanged(c)
			if !changed {
				// No cambió, solo persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			updateData := mapper.MapContactToTable(ctx, updated)
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
