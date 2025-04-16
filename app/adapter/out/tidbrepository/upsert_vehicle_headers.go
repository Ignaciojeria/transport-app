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

type UpsertVehicleHeaders func(context.Context, domain.Headers) (domain.Headers, error)

func init() {
	ioc.Registry(NewUpsertVehicleHeaders, database.NewConnectionFactory)
}
func NewUpsertVehicleHeaders(conn database.ConnectionFactory) UpsertVehicleHeaders {
	return func(ctx context.Context, h domain.Headers) (domain.Headers, error) {
		var vehicleHeaders table.VehicleHeaders
		err := conn.DB.WithContext(ctx).
			Table("vehicle_headers").
			Where(
				"commerce = ? AND consumer = ? AND organization_id = ?",
				h.Commerce, h.Consumer, "TODO").
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
