package usecase

import (
	"context"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type OrderPlannedLog func(ctx context.Context, plan domain.Plan) error

func init() {
	ioc.Registry(NewOrderPlannedLog)
}

func NewOrderPlannedLog() OrderPlannedLog {
	return func(ctx context.Context, plan domain.Plan) error {
		return nil
	}
}
