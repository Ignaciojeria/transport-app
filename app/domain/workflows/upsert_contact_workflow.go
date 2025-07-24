package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertContactWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertContactWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertContactWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertContactWorkflow, error) {
	return UpsertContactWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertContactWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertContactWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertContactStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.ContactUpserted(), Src: []string{w.UpsertContactStarted()}, Dst: w.ContactUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertContactWorkflow) WorkflowName() string {
	return "upsert_contact_workflow"
}

func (w UpsertContactWorkflow) ContactUpserted() string {
	return "contact_upserted"
}

func (w UpsertContactWorkflow) UpsertContactStarted() string {
	return "upsert_contact_started"
}

func (w UpsertContactWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertContactWorkflow) CanTransitionToContactUpserted() bool {
	return w.fsm.Can(w.ContactUpserted())
}

func (w UpsertContactWorkflow) IsContactUpserted() bool {
	return w.fsm.Current() == w.ContactUpserted()
}

func (w UpsertContactWorkflow) IsUpsertContactStarted() bool {
	return w.fsm.Current() == w.UpsertContactStarted()
}

// SetContactUpsertedTransition completa el upsert de contact
func (w UpsertContactWorkflow) SetContactUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.ContactUpserted())
}