package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertPlanHeadersWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertPlanHeadersWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertPlanHeadersWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertPlanHeadersWorkflow, error) {
	return UpsertPlanHeadersWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertPlanHeadersWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertPlanHeadersWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertPlanHeadersStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.PlanHeadersUpserted(), Src: []string{w.UpsertPlanHeadersStarted()}, Dst: w.PlanHeadersUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertPlanHeadersWorkflow) WorkflowName() string {
	return "upsert_plan_headers_workflow"
}

func (w UpsertPlanHeadersWorkflow) PlanHeadersUpserted() string {
	return "plan_headers_upserted"
}

func (w UpsertPlanHeadersWorkflow) UpsertPlanHeadersStarted() string {
	return "upsert_plan_headers_started"
}

func (w UpsertPlanHeadersWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertPlanHeadersWorkflow) CanTransitionToPlanHeadersUpserted() bool {
	return w.fsm.Can(w.PlanHeadersUpserted())
}

func (w UpsertPlanHeadersWorkflow) IsPlanHeadersUpserted() bool {
	return w.fsm.Current() == w.PlanHeadersUpserted()
}

func (w UpsertPlanHeadersWorkflow) IsUpsertPlanHeadersStarted() bool {
	return w.fsm.Current() == w.UpsertPlanHeadersStarted()
}

// SetPlanHeadersUpsertedTransition completa el upsert de plan headers
func (w UpsertPlanHeadersWorkflow) SetPlanHeadersUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.PlanHeadersUpserted())
}