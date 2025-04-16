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

type UpsertVehicleCategory func(context.Context, domain.VehicleCategory) (domain.VehicleCategory, error)

func init() {
	ioc.Registry(NewUpsertVehicleCategory, database.NewConnectionFactory)
}
func NewUpsertVehicleCategory(conn database.ConnectionFactory) UpsertVehicleCategory {
	return func(ctx context.Context, vc domain.VehicleCategory) (domain.VehicleCategory, error) {
		var vehicleCategoryTbl table.VehicleCategory
		err := conn.DB.WithContext(ctx).
			Table("vehicle_categories").
			Where("`type` = ? AND organization_id = ?",
				vc.Type,
				"TODO").
			First(&vehicleCategoryTbl).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.VehicleCategory{}, err
		}
		vehicleCategoryWithChanged, _ := vehicleCategoryTbl.Map().UpdateIfChanged(vc)
		DBVehicleCategoryToUpdate := mapper.MapVehicleCategory(ctx, vehicleCategoryWithChanged)
		DBVehicleCategoryToUpdate.CreatedAt = vehicleCategoryTbl.CreatedAt
		if err := conn.Save(&DBVehicleCategoryToUpdate).Error; err != nil {
			return domain.VehicleCategory{}, err
		}
		return DBVehicleCategoryToUpdate.Map(), nil
	}
}
