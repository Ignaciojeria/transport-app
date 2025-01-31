package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertVehicleCategory func(context.Context, domain.VehicleCategory) (domain.VehicleCategory, error)

func init() {
	ioc.Registry(NewUpsertVehicleCategory, tidb.NewTIDBConnection)
}
func NewUpsertVehicleCategory(conn tidb.TIDBConnection) UpsertVehicleCategory {
	return func(ctx context.Context, vc domain.VehicleCategory) (domain.VehicleCategory, error) {
		var vehicleCategoryTbl table.VehicleCategory
		err := conn.DB.WithContext(ctx).
			Table("vehicle_categories").
			Where("`type` = ? AND organization_country_id = ?",
				vc.Type,
				vc.Organization.OrganizationCountryID).
			First(&vehicleCategoryTbl).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.VehicleCategory{}, err
		}
		vehicleCategoryWithChanged := vehicleCategoryTbl.Map().UpdateIfChanged(vc)
		DBVehicleCategoryToUpdate := mapper.MapVehicleCategory(vehicleCategoryWithChanged)
		DBVehicleCategoryToUpdate.CreatedAt = vehicleCategoryTbl.CreatedAt
		if err := conn.Save(&DBVehicleCategoryToUpdate).Error; err != nil {
			return domain.VehicleCategory{}, err
		}
		return DBVehicleCategoryToUpdate.Map(), nil
	}
}
