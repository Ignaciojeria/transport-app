package usecase

import (
	"context"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertElectricDeliveriesWorkflow func(ctx context.Context, input interface{}) (interface{}, error)

func init() {
	ioc.Registry(NewUpsertElectricDeliveriesWorkflow)
}

func NewUpsertElectricDeliveriesWorkflow() UpsertElectricDeliveriesWorkflow {
	return func(ctx context.Context, input interface{}) (interface{}, error) {
		return input, nil
	}
}
