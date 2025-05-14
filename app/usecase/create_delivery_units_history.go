package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateDeliveryUnitsHistory func(ctx context.Context, input domain.Plan) error

func init() {
	ioc.Registry(NewCreateDeliveryUnitsHistory, tidbrepository.NewUpsertDeliveryUnitsHistory)
}

func NewCreateDeliveryUnitsHistory(saveHistory tidbrepository.UpsertDeliveryUnitsHistory) CreateDeliveryUnitsHistory {
	return func(ctx context.Context, input domain.Plan) error {
		return saveHistory(ctx, input)
	}
}
