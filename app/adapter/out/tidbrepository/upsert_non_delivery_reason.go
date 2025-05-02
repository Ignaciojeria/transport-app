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

type UpsertNonDeliveryReason func(context.Context, domain.NonDeliveryReason) error

func init() {
	ioc.Registry(
		NewUpsertNonDeliveryReason,
		database.NewConnectionFactory,
	)
}

func NewUpsertNonDeliveryReason(conn database.ConnectionFactory) UpsertNonDeliveryReason {
	return func(ctx context.Context, nd domain.NonDeliveryReason) error {
		docID := nd.DocID(ctx)
		if docID.IsZero() {
			return errors.New("cannot persist non delivery reason with empty doc ID")
		}

		var existing table.NonDeliveryReason

		err := conn.DB.WithContext(ctx).
			Table("non_delivery_reasons").
			Where("document_id = ?", docID).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRecord := mapper.MapNonDeliveryReasonTable(ctx, nd)
			return conn.Create(&newRecord).Error
		}

		// Map and compare to check for changes
		existingMapped := existing.Map()
		updated, changed := existingMapped.UpdateIfChanged(nd)

		if !changed {
			return nil
		}

		updateData := mapper.MapNonDeliveryReasonTable(ctx, updated)
		updateData.ID = existing.ID
		updateData.CreatedAt = existing.CreatedAt

		return conn.Save(&updateData).Error
	}
}
