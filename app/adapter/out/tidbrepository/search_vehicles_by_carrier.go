package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchVehiclesByCarrier func(
	context.Context,
	domain.VehicleSearchFilters) ([]domain.Vehicle, error)

func init() {
	ioc.Registry(NewSearchVehiclesByCarrier, tidb.NewTIDBConnection)
}

func NewSearchVehiclesByCarrier(conn tidb.TIDBConnection) SearchVehiclesByCarrier {
	return func(ctx context.Context, vsf domain.VehicleSearchFilters) ([]domain.Vehicle, error) {
		var vehicles []table.Vehicle

		if err := conn.DB.WithContext(ctx).
			Table("vehicles").
			Preload("VehicleCategory").
			Preload("VehicleHeaders").
			Joins("JOIN carriers c ON vehicles.carrier_id = c.id").
			Joins("JOIN organization_countries oc ON vehicles.organization_country_id = oc.id").
			Joins("JOIN organizations org ON oc.organization_id = org.id").
			Joins("JOIN api_keys ak ON org.id = ak.organization_id").
			Where("ak.key = ? AND oc.country = ? AND c.reference_id = ?",
				vsf.Organization.Key,
				vsf.Organization.Country.Alpha2(),
				vsf.CarrierReferenceID).
			Find(&vehicles).Error; err != nil {
			return nil, err
		}

		response := make([]domain.Vehicle, len(vehicles))
		for i, vehicle := range vehicles {
			response[i] = vehicle.Map()
		}

		return response, nil
	}
}
