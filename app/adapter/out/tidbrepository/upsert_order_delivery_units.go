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

type UpsertOrderDeliveryUnits func(context.Context, domain.Order, ...domain.FSMState) error

func init() {
	ioc.Registry(
		NewUpsertOrderDeliveryUnits,
		database.NewConnectionFactory,
		NewSaveFSMTransition)
}
func NewUpsertOrderDeliveryUnits(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertOrderDeliveryUnits {
	return func(ctx context.Context, order domain.Order, fsmState ...domain.FSMState) error {
		orderPackages := mapper.MapOrderDeliveryUnits(ctx, order)
		return conn.Transaction(func(tx *gorm.DB) error {

			if len(orderPackages) > 0 {
				if err := tx.Save(&orderPackages).Error; err != nil {
					return err
				}
			}

			if len(orderPackages) == 0 {
				pkg := domain.DeliveryUnit{}
				if err := tx.Save(&table.OrderDeliveryUnit{
					OrderDoc:        order.DocID(ctx).String(),
					DeliveryUnitDoc: pkg.DocID(ctx).String()}).Error; err != nil {
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
