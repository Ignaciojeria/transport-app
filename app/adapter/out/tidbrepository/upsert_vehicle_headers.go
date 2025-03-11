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

type UpsertVehicleHeaders func(context.Context, domain.Headers) (domain.Headers, error)

func init() {
	ioc.Registry(NewUpsertVehicleHeaders, tidb.NewTIDBConnection)
}
func NewUpsertVehicleHeaders(conn tidb.TIDBConnection) UpsertVehicleHeaders {
	return func(ctx context.Context, h domain.Headers) (domain.Headers, error) {
		var vehicleHeaders table.VehicleHeaders
		err := conn.DB.WithContext(ctx).
			Table("vehicle_headers").
			Where(
				"commerce = ? AND consumer = ? AND organization_id = ?",
				h.Commerce, h.Consumer, h.Organization.ID).
			First(&vehicleHeaders).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Headers{}, err
		}
		if vehicleHeaders.Commerce != "" {
			return vehicleHeaders.Map(), nil
		}
		vehicleHeadersTbl := mapper.MapVehicleHeaders(h)
		if err := conn.Save(&vehicleHeadersTbl).Error; err != nil {
			return domain.Headers{}, err
		}
		return vehicleHeadersTbl.Map(), nil
	}
}
