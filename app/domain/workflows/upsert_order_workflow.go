package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertOrderWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertOrderWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertOrderWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertOrderWorkflow, error) {
	return UpsertOrderWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertOrderWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertOrderWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertOrderStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.OrderUpserted(), Src: []string{w.UpsertOrderStarted()}, Dst: w.OrderUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertOrderWorkflow) WorkflowName() string {
	return "upsert_order_workflow"
}

func (w UpsertOrderWorkflow) OrderUpserted() string {
	return "order_upserted"
}

func (w UpsertOrderWorkflow) UpsertOrderStarted() string {
	return "upsert_order_started"
}

func (w UpsertOrderWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertOrderWorkflow) CanTransitionToOrderUpserted() bool {
	return w.fsm.Can(w.OrderUpserted())
}

func (w UpsertOrderWorkflow) IsOrderUpserted() bool {
	return w.fsm.Current() == w.OrderUpserted()
}

func (w UpsertOrderWorkflow) IsUpsertOrderStarted() bool {
	return w.fsm.Current() == w.UpsertOrderStarted()
}

// SetOrderUpsertedTransition completa el upsert de order
func (w UpsertOrderWorkflow) SetOrderUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OrderUpserted())
}