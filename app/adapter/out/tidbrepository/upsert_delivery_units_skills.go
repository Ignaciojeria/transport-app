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

type UpsertDeliveryUnitsSkills func(context.Context, domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertDeliveryUnitsSkills,
		database.NewConnectionFactory)
}

func NewUpsertDeliveryUnitsSkills(conn database.ConnectionFactory) UpsertDeliveryUnitsSkills {
	return func(ctx context.Context, order domain.Order) error {
		deliveryUnitDocs := make([]string, 0, len(order.DeliveryUnits))
		for _, deliveryUnit := range order.DeliveryUnits {
			deliveryUnitDocs = append(deliveryUnitDocs, deliveryUnit.DocID(ctx).String())
		}

		skills := mapper.MapDeliveryUnitsSkills(ctx, order)
		return conn.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("delivery_unit_doc in ?", deliveryUnitDocs).
				Delete(&table.DeliveryUnitsSkills{}).Error; err != nil {
				return err
			}
			if len(skills) > 0 {
				if err := tx.Save(&skills).Error; err != nil {
					return err
				}
			}
			return nil
		})
	}
}
