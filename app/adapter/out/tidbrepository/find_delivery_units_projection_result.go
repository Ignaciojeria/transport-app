package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/projectionresult"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type FindDeliveryUnitsProjectionResult func(ctx context.Context, filters domain.DeliveryUnitsFilter) (
	projectionresult.DeliveryUnitsProjectionResult, error)

func init() {
	ioc.Registry(
		NewFindDeliveryUnitsProjectionResult,
		database.NewConnectionFactory)
}
func NewFindDeliveryUnitsProjectionResult(
	conn database.ConnectionFactory) FindDeliveryUnitsProjectionResult {
	return func(ctx context.Context, filters domain.DeliveryUnitsFilter) (projectionresult.DeliveryUnitsProjectionResult, error) {

		return projectionresult.DeliveryUnitsProjectionResult{}, nil
	}
}
