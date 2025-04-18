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

type UpsertVehicle func(context.Context, domain.Vehicle) (domain.Vehicle, error)

func init() {
	ioc.Registry(NewUpsertVehicle, database.NewConnectionFactory)
}
func NewUpsertVehicle(conn database.ConnectionFactory) UpsertVehicle {
	return func(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error) {
		vehicle := table.Vehicle{}
		err := conn.DB.WithContext(ctx).Table("vehicles").
			Where("reference_id = ? AND organization_id = ?",
				string("TODO"), "TODO").First(&vehicle).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Vehicle{}, err
		}
		vehicleWithChanges, _ := vehicle.Map().UpdateIfChanged(v)
		DBVehicleToUpsert := mapper.DomainToTableVehicle(vehicleWithChanges)
		DBVehicleToUpsert.CreatedAt = vehicle.CreatedAt
		if err := conn.
			Omit("Organization").
			Omit("VehicleCategory").
			Omit("VehicleHeaders").
			Omit("Carrier").
			Save(&DBVehicleToUpsert).Error; err != nil {
			return domain.Vehicle{}, err
		}
		return DBVehicleToUpsert.Map(), nil
	}
}
