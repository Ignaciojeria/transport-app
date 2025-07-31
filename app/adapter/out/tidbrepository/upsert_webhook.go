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

type UpsertWebhook func(ctx context.Context, webhook domain.Webhook, fsmState ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertWebhook, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertWebhook(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertWebhook {
	return func(ctx context.Context, webhook domain.Webhook, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var existing table.Webhook
			err := tx.WithContext(ctx).
				Table("webhooks").
				Where("document_id = ?", webhook.DocID(ctx)).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// No existe → insert
				newWebhook := mapper.MapWebhookToTable(ctx, webhook)
				if err := tx.Omit("Tenant").Create(&newWebhook).Error; err != nil {
					return err
				}

				// Persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			// Ya existe → update solo si cambió algo
			updated, changed := existing.Map().UpdateIfChanged(webhook)
			if !changed {
				// No hay cambios, solo persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			updateData := mapper.MapWebhookToTable(ctx, updated)
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
