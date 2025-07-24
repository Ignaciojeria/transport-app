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

type UpsertDeliveryUnitsHistory func(ctx context.Context, c domain.Plan, fsmState ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertDeliveryUnitsHistory, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertDeliveryUnitsHistory(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertDeliveryUnitsHistory {
	return func(ctx context.Context, c domain.Plan, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			// Get all delivery units history records for this plan
			deliveryUnitsHistory := mapper.MapDeliveryUnitsHistoryTable(ctx, c)

			// For each delivery unit history record
			for _, duh := range deliveryUnitsHistory {
				var existing table.DeliveryUnitsStatusHistory
				err := tx.WithContext(ctx).
					Table("delivery_units_status_histories").
					Where("document_id = ?", duh.DocumentID).
					First(&existing).Error

				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}

				// No existe → insert
				if err := tx.Create(&duh).Error; err != nil {
					return err
				}
			}

			// Persistir FSMState si está presente
			if len(fsmState) > 0 && saveFSMTransition != nil {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})
	}
}
