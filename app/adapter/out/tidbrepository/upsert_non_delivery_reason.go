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

type UpsertNonDeliveryReason func(context.Context, domain.NonDeliveryReason, ...domain.FSMState) error

func init() {
	ioc.Registry(
		NewUpsertNonDeliveryReason,
		database.NewConnectionFactory,
		NewSaveFSMTransition,
	)
}

func NewUpsertNonDeliveryReason(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertNonDeliveryReason {
	return func(ctx context.Context, nd domain.NonDeliveryReason, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			docID := nd.DocID(ctx)
			if docID.IsZero() {
				return errors.New("cannot persist non delivery reason with empty doc ID")
			}

			var existing table.NonDeliveryReason

			err := tx.WithContext(ctx).
				Table("non_delivery_reasons").
				Where("document_id = ?", docID).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				newRecord := mapper.MapNonDeliveryReasonTable(ctx, nd)
				if err := tx.Create(&newRecord).Error; err != nil {
					return err
				}

				// Persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			// Map and compare to check for changes
			existingMapped := existing.Map()
			updated, changed := existingMapped.UpdateIfChanged(nd)

			if !changed {
				// No hay cambios, solo persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			updateData := mapper.MapNonDeliveryReasonTable(ctx, updated)
			updateData.ID = existing.ID
			updateData.CreatedAt = existing.CreatedAt

			if err := tx.Save(&updateData).Error; err != nil {
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
