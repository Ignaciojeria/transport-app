package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertVehicleWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertVehicleWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertVehicleWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertVehicleWorkflow, error) {
	return UpsertVehicleWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertVehicleWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertVehicleWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertVehicleStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.VehicleUpserted(), Src: []string{w.UpsertVehicleStarted()}, Dst: w.VehicleUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertVehicleWorkflow) WorkflowName() string {
	return "upsert_vehicle_workflow"
}

func (w UpsertVehicleWorkflow) VehicleUpserted() string {
	return "vehicle_upserted"
}

func (w UpsertVehicleWorkflow) UpsertVehicleStarted() string {
	return "upsert_vehicle_started"
}

func (w UpsertVehicleWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertVehicleWorkflow) CanTransitionToVehicleUpserted() bool {
	return w.fsm.Can(w.VehicleUpserted())
}

func (w UpsertVehicleWorkflow) IsVehicleUpserted() bool {
	return w.fsm.Current() == w.VehicleUpserted()
}

func (w UpsertVehicleWorkflow) IsUpsertVehicleStarted() bool {
	return w.fsm.Current() == w.UpsertVehicleStarted()
}

// SetVehicleUpsertedTransition completa el upsert de vehicle
func (w UpsertVehicleWorkflow) SetVehicleUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.VehicleUpserted())
}