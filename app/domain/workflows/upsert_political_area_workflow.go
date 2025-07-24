package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertPoliticalAreaWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertPoliticalAreaWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertPoliticalAreaWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertPoliticalAreaWorkflow, error) {
	return UpsertPoliticalAreaWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertPoliticalAreaWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertPoliticalAreaWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertPoliticalAreaStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.PoliticalAreaUpserted(), Src: []string{w.UpsertPoliticalAreaStarted()}, Dst: w.PoliticalAreaUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertPoliticalAreaWorkflow) WorkflowName() string {
	return "upsert_political_area_workflow"
}

func (w UpsertPoliticalAreaWorkflow) PoliticalAreaUpserted() string {
	return "political_area_upserted"
}

func (w UpsertPoliticalAreaWorkflow) UpsertPoliticalAreaStarted() string {
	return "upsert_political_area_started"
}

func (w UpsertPoliticalAreaWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertPoliticalAreaWorkflow) TransitionToPoliticalAreaUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.PoliticalAreaUpserted())
}

func (w UpsertPoliticalAreaWorkflow) CanTransitionToPoliticalAreaUpserted() bool {
	return w.fsm.Can(w.PoliticalAreaUpserted())
}

func (w UpsertPoliticalAreaWorkflow) IsPoliticalAreaUpserted() bool {
	return w.fsm.Current() == w.PoliticalAreaUpserted()
}

func (w UpsertPoliticalAreaWorkflow) IsUpsertPoliticalAreaStarted() bool {
	return w.fsm.Current() == w.UpsertPoliticalAreaStarted()
}

// SetPoliticalAreaUpsertedTransition completa el upsert de political area
func (w UpsertPoliticalAreaWorkflow) SetPoliticalAreaUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.PoliticalAreaUpserted())
}