package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertSkillWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertSkillWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertSkillWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertSkillWorkflow, error) {
	return UpsertSkillWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertSkillWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertSkillWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertSkillStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.SkillUpserted(), Src: []string{w.UpsertSkillStarted()}, Dst: w.SkillUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertSkillWorkflow) WorkflowName() string {
	return "upsert_skill_workflow"
}

func (w UpsertSkillWorkflow) SkillUpserted() string {
	return "skill_upserted"
}

func (w UpsertSkillWorkflow) UpsertSkillStarted() string {
	return "upsert_skill_started"
}

func (w UpsertSkillWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertSkillWorkflow) TransitionToSkillUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.SkillUpserted())
}

func (w UpsertSkillWorkflow) CanTransitionToSkillUpserted() bool {
	return w.fsm.Can(w.SkillUpserted())
}

func (w UpsertSkillWorkflow) IsSkillUpserted() bool {
	return w.fsm.Current() == w.SkillUpserted()
}

func (w UpsertSkillWorkflow) IsUpsertSkillStarted() bool {
	return w.fsm.Current() == w.UpsertSkillStarted()
}

// SetSkillUpsertedTransition completa el upsert de skill
func (w UpsertSkillWorkflow) SetSkillUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.SkillUpserted())
}