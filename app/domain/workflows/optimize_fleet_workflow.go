package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type OptimizeFleetWorkflow struct {
	IdempotencyKey       string
	NextInput            []byte
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewOptimizeFleetWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewOptimizeFleetWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (OptimizeFleetWorkflow, error) {
	return OptimizeFleetWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w OptimizeFleetWorkflow) Restore(ctx context.Context, idempotencyKey string) (OptimizeFleetWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.OptimizationStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.OptimizationCompleted(), Src: []string{w.OptimizationStarted()}, Dst: w.OptimizationCompleted()},
		},
		fsm.Callbacks{},
	)
	w.NextInput = lastTransition.NextInput
	return w, nil
}

func (w OptimizeFleetWorkflow) WorkflowName() string {
	return "optimize_fleet_workflow"
}

func (w OptimizeFleetWorkflow) OptimizationCompleted() string {
	return "optimization_completed"
}

func (w OptimizeFleetWorkflow) OptimizationStarted() string {
	return "optimization_started"
}

func (w OptimizeFleetWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		NextInput:      w.NextInput,
		State:          w.fsm.Current(),
	}
}

func (w OptimizeFleetWorkflow) TransitionToOptimizationCompleted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OptimizationCompleted())
}

func (w OptimizeFleetWorkflow) CanTransitionToOptimizationCompleted() bool {
	return w.fsm.Can(w.OptimizationCompleted())
}

func (w OptimizeFleetWorkflow) IsOptimizationCompleted() bool {
	return w.fsm.Current() == w.OptimizationCompleted()
}

func (w OptimizeFleetWorkflow) IsOptimizationStarted() bool {
	return w.fsm.Current() == w.OptimizationStarted()
}

// CompleteOptimization completa la optimización de la flota
func (w OptimizeFleetWorkflow) SetOptimizationCompletedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.OptimizationCompleted())
}
