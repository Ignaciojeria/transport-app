package order

import (
	"context"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CancelOrder func(ctx context.Context, aggregateID int64, request events.OrderCancelledRequest) error

func init() {
	ioc.Registry(NewCancelOrder,
		observability.NewObservability,
		supabaserepo.NewUpdateOrderStatus,
	)
}

func NewCancelOrder(
	observability observability.Observability,
	updateOrderStatus supabaserepo.UpdateOrderStatus) CancelOrder {
	return func(ctx context.Context, aggregateID int64, request events.OrderCancelledRequest) error {
		observability.Logger.InfoContext(ctx, "cancel_order", "aggregateID", aggregateID, "request", request)
		spanCtx, span := observability.Tracer.Start(ctx, "cancel_order")
		defer span.End()

		// Para cancelar, actualizamos todos los items que no estén en DISPATCHED o CANCELLED
		err := updateOrderStatus(
			spanCtx,
			aggregateID,
			"CANCELLED",
			nil, // nil significa todos los items
			"",  // No filtramos por estación
			events.EventOrderCancelled,
			request,
		)
		if err != nil {
			observability.Logger.ErrorContext(spanCtx, "error cancelling order", "error", err)
			return err
		}

		observability.Logger.InfoContext(spanCtx, "order cancelled successfully", "aggregateID", aggregateID)
		return nil
	}
}
