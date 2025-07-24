package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertVehicleCategoryWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertVehicleCategoryWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertVehicleCategoryWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertVehicleCategoryWorkflow, error) {
	return UpsertVehicleCategoryWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertVehicleCategoryWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertVehicleCategoryWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertVehicleCategoryStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.VehicleCategoryUpserted(), Src: []string{w.UpsertVehicleCategoryStarted()}, Dst: w.VehicleCategoryUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertVehicleCategoryWorkflow) WorkflowName() string {
	return "upsert_vehicle_category_workflow"
}

func (w UpsertVehicleCategoryWorkflow) VehicleCategoryUpserted() string {
	return "vehicle_category_upserted"
}

func (w UpsertVehicleCategoryWorkflow) UpsertVehicleCategoryStarted() string {
	return "upsert_vehicle_category_started"
}

func (w UpsertVehicleCategoryWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertVehicleCategoryWorkflow) TransitionToVehicleCategoryUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.VehicleCategoryUpserted())
}

func (w UpsertVehicleCategoryWorkflow) CanTransitionToVehicleCategoryUpserted() bool {
	return w.fsm.Can(w.VehicleCategoryUpserted())
}

func (w UpsertVehicleCategoryWorkflow) IsVehicleCategoryUpserted() bool {
	return w.fsm.Current() == w.VehicleCategoryUpserted()
}

func (w UpsertVehicleCategoryWorkflow) IsUpsertVehicleCategoryStarted() bool {
	return w.fsm.Current() == w.UpsertVehicleCategoryStarted()
}

// SetVehicleCategoryUpsertedTransition completa el upsert de vehicle category
func (w UpsertVehicleCategoryWorkflow) SetVehicleCategoryUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.VehicleCategoryUpserted())
}