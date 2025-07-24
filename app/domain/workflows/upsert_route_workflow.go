package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertRouteWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertRouteWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertRouteWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertRouteWorkflow, error) {
	return UpsertRouteWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertRouteWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertRouteWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertRouteStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.RouteUpserted(), Src: []string{w.UpsertRouteStarted()}, Dst: w.RouteUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertRouteWorkflow) WorkflowName() string {
	return "upsert_route_workflow"
}

func (w UpsertRouteWorkflow) RouteUpserted() string {
	return "route_upserted"
}

func (w UpsertRouteWorkflow) UpsertRouteStarted() string {
	return "upsert_route_started"
}

func (w UpsertRouteWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertRouteWorkflow) CanTransitionToRouteUpserted() bool {
	return w.fsm.Can(w.RouteUpserted())
}

func (w UpsertRouteWorkflow) IsRouteUpserted() bool {
	return w.fsm.Current() == w.RouteUpserted()
}

func (w UpsertRouteWorkflow) IsUpsertRouteStarted() bool {
	return w.fsm.Current() == w.UpsertRouteStarted()
}

// SetRouteUpsertedTransition completa el upsert de route
func (w UpsertRouteWorkflow) SetRouteUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.RouteUpserted())
}