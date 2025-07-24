package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertOrderReferencesWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertOrderReferencesWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertOrderReferencesWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertOrderReferencesWorkflow, error) {
	return UpsertOrderReferencesWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertOrderReferencesWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertOrderReferencesWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertOrderReferencesStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.OrderReferencesUpserted(), Src: []string{w.UpsertOrderReferencesStarted()}, Dst: w.OrderReferencesUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertOrderReferencesWorkflow) WorkflowName() string {
	return "upsert_order_references_workflow"
}

func (w UpsertOrderReferencesWorkflow) OrderReferencesUpserted() string {
	return "order_references_upserted"
}

func (w UpsertOrderReferencesWorkflow) UpsertOrderReferencesStarted() string {
	return "upsert_order_references_started"
}

func (w UpsertOrderReferencesWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertOrderReferencesWorkflow) CanTransitionToOrderReferencesUpserted() bool {
	return w.fsm.Can(w.OrderReferencesUpserted())
}

func (w UpsertOrderReferencesWorkflow) IsOrderReferencesUpserted() bool {
	return w.fsm.Current() == w.OrderReferencesUpserted()
}

func (w UpsertOrderReferencesWorkflow) IsUpsertOrderReferencesStarted() bool {
	return w.fsm.Current() == w.UpsertOrderReferencesStarted()
}

// SetOrderReferencesUpsertedTransition completa el upsert de order references
func (w UpsertOrderReferencesWorkflow) SetOrderReferencesUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OrderReferencesUpserted())
}