package order

import (
	"context"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateOrderResult struct {
	OrderNumber int `json:"orderNumber"`
}

type CreateOrder func(ctx context.Context, menuID string, request events.CreateOrderRequest) (CreateOrderResult, error)

func init() {
	ioc.Registry(NewCreateOrder,
		observability.NewObservability,
		supabaserepo.NewSaveMenuOrder,
	)
}

func NewCreateOrder(
	observability observability.Observability,
	saveMenuOrder supabaserepo.SaveMenuOrder) CreateOrder {
	return func(ctx context.Context, menuID string, request events.CreateOrderRequest) (CreateOrderResult, error) {
		observability.Logger.InfoContext(ctx, "create_order", "menuID", menuID, "request", request)
		spanCtx, span := observability.Tracer.Start(ctx, "create_order")
		defer span.End()

		result, err := saveMenuOrder(
			spanCtx,
			menuID,
			request,
			events.EventCreateOrderRequested,
		)
		if err != nil {
			observability.Logger.ErrorContext(spanCtx, "error creating order", "error", err)
			return CreateOrderResult{}, err
		}

		observability.Logger.InfoContext(spanCtx, "order created successfully", "menuID", menuID, "orderNumber", result.OrderNumber)
		return CreateOrderResult{
			OrderNumber: result.OrderNumber,
		}, nil
	}
}
