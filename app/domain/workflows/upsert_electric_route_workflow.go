package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertElectricRouteWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertElectricRouteWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertElectricRouteWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertElectricRouteWorkflow, error) {
	return UpsertElectricRouteWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertElectricRouteWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertElectricRouteWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertElectricRouteStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.ElectricRouteUpserted(), Src: []string{w.UpsertElectricRouteStarted()}, Dst: w.ElectricRouteUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertElectricRouteWorkflow) WorkflowName() string {
	return "upsert_electric_route_workflow"
}

func (w UpsertElectricRouteWorkflow) ElectricRouteUpserted() string {
	return "electric_route_upserted"
}

func (w UpsertElectricRouteWorkflow) UpsertElectricRouteStarted() string {
	return "upsert_electric_route_started"
}

func (w UpsertElectricRouteWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertElectricRouteWorkflow) CanTransitionToElectricRouteUpserted() bool {
	return w.fsm.Can(w.ElectricRouteUpserted())
}

func (w UpsertElectricRouteWorkflow) IsElectricRouteUpserted() bool {
	return w.fsm.Current() == w.ElectricRouteUpserted()
}

func (w UpsertElectricRouteWorkflow) IsUpsertElectricRouteStarted() bool {
	return w.fsm.Current() == w.UpsertElectricRouteStarted()
}

// SetElectricRouteUpsertedTransition completa el upsert de electric route
func (w UpsertElectricRouteWorkflow) SetElectricRouteUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.ElectricRouteUpserted())
}
