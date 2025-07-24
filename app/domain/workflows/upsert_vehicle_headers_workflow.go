package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertVehicleHeadersWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertVehicleHeadersWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertVehicleHeadersWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertVehicleHeadersWorkflow, error) {
	return UpsertVehicleHeadersWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertVehicleHeadersWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertVehicleHeadersWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertVehicleHeadersStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.VehicleHeadersUpserted(), Src: []string{w.UpsertVehicleHeadersStarted()}, Dst: w.VehicleHeadersUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertVehicleHeadersWorkflow) WorkflowName() string {
	return "upsert_vehicle_headers_workflow"
}

func (w UpsertVehicleHeadersWorkflow) VehicleHeadersUpserted() string {
	return "vehicle_headers_upserted"
}

func (w UpsertVehicleHeadersWorkflow) UpsertVehicleHeadersStarted() string {
	return "upsert_vehicle_headers_started"
}

func (w UpsertVehicleHeadersWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertVehicleHeadersWorkflow) TransitionToVehicleHeadersUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.VehicleHeadersUpserted())
}

func (w UpsertVehicleHeadersWorkflow) CanTransitionToVehicleHeadersUpserted() bool {
	return w.fsm.Can(w.VehicleHeadersUpserted())
}

func (w UpsertVehicleHeadersWorkflow) IsVehicleHeadersUpserted() bool {
	return w.fsm.Current() == w.VehicleHeadersUpserted()
}

func (w UpsertVehicleHeadersWorkflow) IsUpsertVehicleHeadersStarted() bool {
	return w.fsm.Current() == w.UpsertVehicleHeadersStarted()
}

// SetVehicleHeadersUpsertedTransition completa el upsert de vehicle headers
func (w UpsertVehicleHeadersWorkflow) SetVehicleHeadersUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.VehicleHeadersUpserted())
}