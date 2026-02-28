package order

import (
	"context"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/ioc"
)

type MarkReady func(ctx context.Context, aggregateID int64, request events.OrderItemReadyRequest) error

func init() {
	ioc.Register(NewMarkReady)
}

func NewMarkReady(
	observability observability.Observability,
	updateOrderStatus supabaserepo.UpdateOrderStatus) MarkReady {
	return func(ctx context.Context, aggregateID int64, request events.OrderItemReadyRequest) error {
		observability.Logger.InfoContext(ctx, "mark_ready", "aggregateID", aggregateID, "request", request)
		spanCtx, span := observability.Tracer.Start(ctx, "mark_ready")
		defer span.End()

		err := updateOrderStatus(
			spanCtx,
			aggregateID,
			"READY",
			request.ItemKeys,
			request.Station,
			events.EventOrderItemReady,
			request,
		)
		if err != nil {
			observability.Logger.ErrorContext(spanCtx, "error marking ready", "error", err)
			return err
		}

		observability.Logger.InfoContext(spanCtx, "items marked ready successfully", "aggregateID", aggregateID)
		return nil
	}
}
