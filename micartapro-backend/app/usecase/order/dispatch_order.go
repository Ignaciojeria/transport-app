package order

import (
	"context"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type DispatchOrder func(ctx context.Context, aggregateID int64, request events.OrderDispatchedRequest) error

func init() {
	ioc.Registry(NewDispatchOrder,
		observability.NewObservability,
		supabaserepo.NewUpdateOrderStatus,
	)
}

func NewDispatchOrder(
	observability observability.Observability,
	updateOrderStatus supabaserepo.UpdateOrderStatus) DispatchOrder {
	return func(ctx context.Context, aggregateID int64, request events.OrderDispatchedRequest) error {
		observability.Logger.InfoContext(ctx, "dispatch_order", "aggregateID", aggregateID, "request", request)
		spanCtx, span := observability.Tracer.Start(ctx, "dispatch_order")
		defer span.End()

		// Para dispatch, actualizamos todos los items que estén en READY a DISPATCHED
		err := updateOrderStatus(
			spanCtx,
			aggregateID,
			"DISPATCHED",
			nil, // nil significa todos los items
			"",  // No filtramos por estación
			events.EventOrderDispatched,
			request,
		)
		if err != nil {
			observability.Logger.ErrorContext(spanCtx, "error dispatching order", "error", err)
			return err
		}

		observability.Logger.InfoContext(spanCtx, "order dispatched successfully", "aggregateID", aggregateID)
		return nil
	}
}
