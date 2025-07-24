package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertDeliveryUnitsLabels func(context.Context, domain.Order, ...domain.FSMState) error

func init() {
	ioc.Registry(
		NewUpsertDeliveryUnitsLabels,
		database.NewConnectionFactory,
		NewSaveFSMTransition)
}

func NewUpsertDeliveryUnitsLabels(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertDeliveryUnitsLabels {
	return func(ctx context.Context, order domain.Order, fsmState ...domain.FSMState) error {
		deliveryUnitDocs := make([]string, 0, len(order.DeliveryUnits))
		for _, deliveryUnit := range order.DeliveryUnits {
			deliveryUnitDocs = append(deliveryUnitDocs, deliveryUnit.DocID(ctx).String())
		}

		labels := mapper.MapDeliveryUnitsLabels(ctx, order)
		return conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("delivery_unit_doc in ?", deliveryUnitDocs).
				Delete(&table.DeliveryUnitsLabels{}).Error; err != nil {
				return err
			}
			if len(labels) > 0 {
				if err := tx.Save(&labels).Error; err != nil {
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
