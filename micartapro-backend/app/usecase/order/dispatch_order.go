package order

import (
	"context"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/ioc"
)

type DispatchOrder func(ctx context.Context, aggregateID int64, request events.OrderDeliveredRequest) error

func init() {
	ioc.Register(NewDispatchOrder)
}

func NewDispatchOrder(
	observability observability.Observability,
	updateOrderStatus supabaserepo.UpdateOrderStatus) DispatchOrder {
	return func(ctx context.Context, aggregateID int64, request events.OrderDeliveredRequest) error {
		observability.Logger.InfoContext(ctx, "dispatch_order", "aggregateID", aggregateID, "request", request)
		spanCtx, span := observability.Tracer.Start(ctx, "dispatch_order")
		defer span.End()

		// COMPLETE: la DB pone DELIVERED (PICKUP/retiro, DIGITAL/productos digitales) o DISPATCHED (DELIVERY/despacho) según fulfillment
		err := updateOrderStatus(
			spanCtx,
			aggregateID,
			"COMPLETE",
			nil, // nil significa todos los items
			"",  // No filtramos por estación
			events.EventOrderDelivered,
			request,
		)
		if err != nil {
			observability.Logger.ErrorContext(spanCtx, "error dispatching order", "error", err)
			return err
		}

		observability.Logger.InfoContext(spanCtx, "order delivered successfully", "aggregateID", aggregateID)
		return nil
	}
}
