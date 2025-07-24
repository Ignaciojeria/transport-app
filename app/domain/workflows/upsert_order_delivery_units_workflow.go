package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertOrderDeliveryUnitsWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertOrderDeliveryUnitsWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertOrderDeliveryUnitsWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertOrderDeliveryUnitsWorkflow, error) {
	return UpsertOrderDeliveryUnitsWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertOrderDeliveryUnitsWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertOrderDeliveryUnitsWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertOrderDeliveryUnitsStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.OrderDeliveryUnitsUpserted(), Src: []string{w.UpsertOrderDeliveryUnitsStarted()}, Dst: w.OrderDeliveryUnitsUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertOrderDeliveryUnitsWorkflow) WorkflowName() string {
	return "upsert_order_delivery_units_workflow"
}

func (w UpsertOrderDeliveryUnitsWorkflow) OrderDeliveryUnitsUpserted() string {
	return "order_delivery_units_upserted"
}

func (w UpsertOrderDeliveryUnitsWorkflow) UpsertOrderDeliveryUnitsStarted() string {
	return "upsert_order_delivery_units_started"
}

func (w UpsertOrderDeliveryUnitsWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertOrderDeliveryUnitsWorkflow) TransitionToOrderDeliveryUnitsUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OrderDeliveryUnitsUpserted())
}

func (w UpsertOrderDeliveryUnitsWorkflow) CanTransitionToOrderDeliveryUnitsUpserted() bool {
	return w.fsm.Can(w.OrderDeliveryUnitsUpserted())
}

func (w UpsertOrderDeliveryUnitsWorkflow) IsOrderDeliveryUnitsUpserted() bool {
	return w.fsm.Current() == w.OrderDeliveryUnitsUpserted()
}

func (w UpsertOrderDeliveryUnitsWorkflow) IsUpsertOrderDeliveryUnitsStarted() bool {
	return w.fsm.Current() == w.UpsertOrderDeliveryUnitsStarted()
}

// SetOrderDeliveryUnitsUpsertedTransition completa el upsert de order delivery units
func (w UpsertOrderDeliveryUnitsWorkflow) SetOrderDeliveryUnitsUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OrderDeliveryUnitsUpserted())
}