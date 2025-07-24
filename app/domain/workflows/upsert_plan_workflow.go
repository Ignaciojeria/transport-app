package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertPlanWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertPlanWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertPlanWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertPlanWorkflow, error) {
	return UpsertPlanWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertPlanWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertPlanWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertPlanStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.PlanUpserted(), Src: []string{w.UpsertPlanStarted()}, Dst: w.PlanUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertPlanWorkflow) WorkflowName() string {
	return "upsert_plan_workflow"
}

func (w UpsertPlanWorkflow) PlanUpserted() string {
	return "plan_upserted"
}

func (w UpsertPlanWorkflow) UpsertPlanStarted() string {
	return "upsert_plan_started"
}

func (w UpsertPlanWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertPlanWorkflow) TransitionToPlanUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.PlanUpserted())
}

func (w UpsertPlanWorkflow) CanTransitionToPlanUpserted() bool {
	return w.fsm.Can(w.PlanUpserted())
}

func (w UpsertPlanWorkflow) IsPlanUpserted() bool {
	return w.fsm.Current() == w.PlanUpserted()
}

func (w UpsertPlanWorkflow) IsUpsertPlanStarted() bool {
	return w.fsm.Current() == w.UpsertPlanStarted()
}

// SetPlanUpsertedTransition completa el upsert de plan
func (w UpsertPlanWorkflow) SetPlanUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.PlanUpserted())
}