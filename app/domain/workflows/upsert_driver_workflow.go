package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertDriverWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertDriverWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertDriverWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertDriverWorkflow, error) {
	return UpsertDriverWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertDriverWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertDriverWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertDriverStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.DriverUpserted(), Src: []string{w.UpsertDriverStarted()}, Dst: w.DriverUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertDriverWorkflow) WorkflowName() string {
	return "upsert_driver_workflow"
}

func (w UpsertDriverWorkflow) DriverUpserted() string {
	return "driver_upserted"
}

func (w UpsertDriverWorkflow) UpsertDriverStarted() string {
	return "upsert_driver_started"
}

func (w UpsertDriverWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertDriverWorkflow) TransitionToDriverUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.DriverUpserted())
}

func (w UpsertDriverWorkflow) CanTransitionToDriverUpserted() bool {
	return w.fsm.Can(w.DriverUpserted())
}

func (w UpsertDriverWorkflow) IsDriverUpserted() bool {
	return w.fsm.Current() == w.DriverUpserted()
}

func (w UpsertDriverWorkflow) IsUpsertDriverStarted() bool {
	return w.fsm.Current() == w.UpsertDriverStarted()
}

// SetDriverUpsertedTransition completa el upsert de driver
func (w UpsertDriverWorkflow) SetDriverUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.DriverUpserted())
}