package order

import (
	"context"
	"errors"
	"micartapro/app/adapter/out/supabaserepo"
	"micartapro/app/events"
	"micartapro/app/shared/infrastructure/observability"
	"micartapro/app/shared/tracking"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

const maxTrackingIDRetries = 5

type CreateOrderResult struct {
	OrderNumber int64  `json:"orderNumber"`
	AggregateID int64  `json:"aggregateId"`
	TrackingID  string `json:"trackingId"`
}

type CreateOrder func(ctx context.Context, menuID string, request events.CreateOrderRequest) (CreateOrderResult, error)

func init() {
	ioc.Registry(NewCreateOrder,
		observability.NewObservability,
		supabaserepo.NewSaveMenuOrder,
		supabaserepo.NewInsertOrderTracking,
	)
}

func NewCreateOrder(
	observability observability.Observability,
	saveMenuOrder supabaserepo.SaveMenuOrder,
	insertOrderTracking supabaserepo.InsertOrderTracking,
) CreateOrder {
	return func(ctx context.Context, menuID string, request events.CreateOrderRequest) (CreateOrderResult, error) {
		observability.Logger.InfoContext(ctx, "create_order", "menuID", menuID, "request", request)
		spanCtx, span := observability.Tracer.Start(ctx, "create_order")
		defer span.End()

		// 1) Evento OrderCreated primero. No contiene tracking_id (la tabla order_tracking es proyección de lectura).
		payload := request
		payload.TrackingID = ""
		result, err := saveMenuOrder(
			spanCtx,
			menuID,
			payload,
			events.EventCreateOrderRequested,
		)
		if err != nil {
			observability.Logger.ErrorContext(spanCtx, "error creating order", "error", err)
			return CreateOrderResult{}, err
		}

		// 2) Side-effect síncrono del command handler: poblar order_tracking para acceso por token.
		// No es un proyector async (eso sería outbox/subscriber, reproyectable). Si el proceso cae entre
		// commit del evento y este insert, la orden existe pero no hay fila en order_tracking; un backfill
		// o GET /track/{id} que regenere bajo demanda cierra el hueco (MVP).
		// Retry si colisión en tracking_id (UNIQUE).
		var trackingID string
		for attempt := 0; attempt < maxTrackingIDRetries; attempt++ {
			tid, err := tracking.GenerateTrackingID()
			if err != nil {
				observability.Logger.ErrorContext(spanCtx, "error generating trackingId", "error", err)
				return CreateOrderResult{}, err
			}
			err = insertOrderTracking(spanCtx, result.AggregateID, tid)
			if err == nil {
				trackingID = tid
				break
			}
			if errors.Is(err, supabaserepo.ErrTrackingIDConflict) {
				observability.Logger.WarnContext(spanCtx, "tracking_id collision, retrying", "attempt", attempt+1)
				continue
			}
			observability.Logger.ErrorContext(spanCtx, "error inserting order_tracking", "error", err)
			return CreateOrderResult{}, err
		}
		if trackingID == "" {
			observability.Logger.ErrorContext(spanCtx, "could not generate unique tracking_id after retries")
			return CreateOrderResult{}, errors.New("could not generate unique tracking_id")
		}

		observability.Logger.InfoContext(spanCtx, "order created successfully", "menuID", menuID, "orderNumber", result.OrderNumber, "aggregateID", result.AggregateID, "trackingID", trackingID)
		return CreateOrderResult{
			OrderNumber: result.AggregateID,
			AggregateID: result.AggregateID,
			TrackingID:  trackingID,
		}, nil
	}
}
