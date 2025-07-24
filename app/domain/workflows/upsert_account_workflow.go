package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertAccountWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertAccountWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertAccountWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertAccountWorkflow, error) {
	return UpsertAccountWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertAccountWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertAccountWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertAccountStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.AccountUpserted(), Src: []string{w.UpsertAccountStarted()}, Dst: w.AccountUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertAccountWorkflow) WorkflowName() string {
	return "upsert_account_workflow"
}

func (w UpsertAccountWorkflow) AccountUpserted() string {
	return "account_upserted"
}

func (w UpsertAccountWorkflow) UpsertAccountStarted() string {
	return "upsert_account_started"
}

func (w UpsertAccountWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertAccountWorkflow) TransitionToAccountUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.AccountUpserted())
}

func (w UpsertAccountWorkflow) CanTransitionToAccountUpserted() bool {
	return w.fsm.Can(w.AccountUpserted())
}

func (w UpsertAccountWorkflow) IsAccountUpserted() bool {
	return w.fsm.Current() == w.AccountUpserted()
}

func (w UpsertAccountWorkflow) IsUpsertAccountStarted() bool {
	return w.fsm.Current() == w.UpsertAccountStarted()
}

// SetAccountUpsertedTransition completa el upsert de account
func (w UpsertAccountWorkflow) SetAccountUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.AccountUpserted())
}