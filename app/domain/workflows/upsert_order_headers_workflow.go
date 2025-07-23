package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertOrderHeadersWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertOrderHeadersWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertOrderHeadersWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertOrderHeadersWorkflow, error) {
	return UpsertOrderHeadersWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertOrderHeadersWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertOrderHeadersWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertOrderHeadersStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.OrderHeadersUpserted(), Src: []string{w.UpsertOrderHeadersStarted()}, Dst: w.OrderHeadersUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertOrderHeadersWorkflow) WorkflowName() string {
	return "upsert_order_headers_workflow"
}

func (w UpsertOrderHeadersWorkflow) OrderHeadersUpserted() string {
	return "order_headers_upserted"
}

func (w UpsertOrderHeadersWorkflow) UpsertOrderHeadersStarted() string {
	return "upsert_order_headers_started"
}

func (w UpsertOrderHeadersWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertOrderHeadersWorkflow) TransitionToOrderHeadersUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OrderHeadersUpserted())
}

func (w UpsertOrderHeadersWorkflow) CanTransitionToOrderHeadersUpserted() bool {
	return w.fsm.Can(w.OrderHeadersUpserted())
}

func (w UpsertOrderHeadersWorkflow) IsOrderHeadersUpserted() bool {
	return w.fsm.Current() == w.OrderHeadersUpserted()
}

func (w UpsertOrderHeadersWorkflow) IsUpsertOrderHeadersStarted() bool {
	return w.fsm.Current() == w.UpsertOrderHeadersStarted()
}

// SetOrderHeadersUpsertedTransition completa el upsert de order headers
func (w UpsertOrderHeadersWorkflow) SetOrderHeadersUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OrderHeadersUpserted())
}
