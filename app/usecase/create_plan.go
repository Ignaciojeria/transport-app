package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreatePlan func(ctx context.Context, input domain.Plan) error

func init() {
	ioc.Registry(
		NewCreatePlan,
		tidbrepository.NewUpsertDeliveryUnitsHistory)
}

func NewCreatePlan(
	upsertDeliveryUnitsHistory tidbrepository.UpsertDeliveryUnitsHistory) CreatePlan {
	return func(ctx context.Context, input domain.Plan) error {
		input.AssignIndexesToAllOrders()
		input.AssignSequenceNumbersToAllOrders()
		err := upsertDeliveryUnitsHistory(ctx, input)
		if err != nil {
			return err
		}
		fmt.Println("planID:" + input.ReferenceID)
		return nil
	}
}
