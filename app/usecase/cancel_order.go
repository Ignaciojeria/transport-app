package usecase

import (
	"context"
	"fmt"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CancelOrder func(ctx context.Context, input domain.Route) error

func init() {
	ioc.Registry(NewCancelOrder)
}

func NewCancelOrder() CancelOrder {
	return func(ctx context.Context, input domain.Route) error {
		fmt.Print("works")
		return nil
	}
}
