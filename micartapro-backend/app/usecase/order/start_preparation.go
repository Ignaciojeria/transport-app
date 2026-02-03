package order

import (
	"context"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type StartPreparation func(ctx context.Context, aggregateID int64, request events.OrderStartedPreparationRequest) error

func init() {
	ioc.Registry(NewStartPreparation,
		observability.NewObservability,
		supabaserepo.NewUpdateOrderStatus,
	)
}

func NewStartPreparation(
	observability observability.Observability,
	updateOrderStatus supabaserepo.UpdateOrderStatus) StartPreparation {
	return func(ctx context.Context, aggregateID int64, request events.OrderStartedPreparationRequest) error {
		observability.Logger.InfoContext(ctx, "start_preparation", "aggregateID", aggregateID, "request", request)
		spanCtx, span := observability.Tracer.Start(ctx, "start_preparation")
		defer span.End()

		err := updateOrderStatus(
			spanCtx,
			aggregateID,
			"IN_PROGRESS",
			request.ItemKeys,
			request.Station,
			events.EventOrderStartedPreparation,
			request,
		)
		if err != nil {
			observability.Logger.ErrorContext(spanCtx, "error starting preparation", "error", err)
			return err
		}

		observability.Logger.InfoContext(spanCtx, "preparation started successfully", "aggregateID", aggregateID)
		return nil
	}
}
