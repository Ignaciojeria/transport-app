package usecase

import (
	"context"
	"fmt"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type ConfirmDeliveries func(ctx context.Context, input domain.Route) error

func init() {
	ioc.Registry(NewConfirmDeliveries)
}

func NewConfirmDeliveries() ConfirmDeliveries {
	return func(ctx context.Context, input domain.Route) error {
		fmt.Println(input)
		return nil
	}
}
