package tidbrepository

import (
	"context"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertVehicleQuery func(context.Context, domain.Vehicle) (domain.Vehicle, error)

func init() {
	ioc.Registry(NewUpsertVehicleQuery)
}
func NewUpsertVehicleQuery() UpsertVehicleQuery {
	return func(ctx context.Context, v domain.Vehicle) (domain.Vehicle, error) {
		return domain.Vehicle{}, nil
	}
}
