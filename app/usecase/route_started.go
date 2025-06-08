package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type RouteStarted func(ctx context.Context, input domain.Route) error

func init() {
	ioc.Registry(NewRouteStarted, tidbrepository.NewUpsertDeliveryUnitsHistory)
}

func NewRouteStarted(
	upsertDeliveryUnitsHistory tidbrepository.UpsertDeliveryUnitsHistory) RouteStarted {
	return func(ctx context.Context, input domain.Route) error {
		fmt.Println("works")
		return nil
	}
}
