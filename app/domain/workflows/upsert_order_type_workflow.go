package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertOrderTypeWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertOrderTypeWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertOrderTypeWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertOrderTypeWorkflow, error) {
	return UpsertOrderTypeWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertOrderTypeWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertOrderTypeWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertOrderTypeStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.OrderTypeUpserted(), Src: []string{w.UpsertOrderTypeStarted()}, Dst: w.OrderTypeUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertOrderTypeWorkflow) WorkflowName() string {
	return "upsert_order_type_workflow"
}

func (w UpsertOrderTypeWorkflow) OrderTypeUpserted() string {
	return "order_type_upserted"
}

func (w UpsertOrderTypeWorkflow) UpsertOrderTypeStarted() string {
	return "upsert_order_type_started"
}

func (w UpsertOrderTypeWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertOrderTypeWorkflow) TransitionToOrderTypeUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OrderTypeUpserted())
}

func (w UpsertOrderTypeWorkflow) CanTransitionToOrderTypeUpserted() bool {
	return w.fsm.Can(w.OrderTypeUpserted())
}

func (w UpsertOrderTypeWorkflow) IsOrderTypeUpserted() bool {
	return w.fsm.Current() == w.OrderTypeUpserted()
}

func (w UpsertOrderTypeWorkflow) IsUpsertOrderTypeStarted() bool {
	return w.fsm.Current() == w.UpsertOrderTypeStarted()
}

// SetOrderTypeUpsertedTransition completa el upsert de order type
func (w UpsertOrderTypeWorkflow) SetOrderTypeUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OrderTypeUpserted())
}