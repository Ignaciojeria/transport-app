package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertVehicleQuery func(context.Context, domain.Vehicle) (domain.Vehicle, error)

func init() {
	ioc.Registry(NewUpsertVehicleQuery, tidb.NewTIDBConnection)
}

func NewUpsertVehicleQuery(conn tidb.TIDBConnection) UpsertVehicleQuery {
	return func(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error) {
		vehicle := mapper.DomainToTableVehicle(v)

		// Buscar vehículo por placa
		if err := conn.WithContext(ctx).
			Where("plate = ?", v.Plate).
			First(&vehicle).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.Vehicle{}, err
			}
		}

		// Buscar carrier por nationalID
		if err := conn.WithContext(ctx).
			Where("national_id = ?", v.Carrier.NationalID).
			First(&vehicle.Carrier).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.Vehicle{}, err
			}
		}

		// Mapear de vuelta a dominio
		return mapper.TableToDomainVehicle(vehicle, vehicle.Carrier), nil
	}
}
