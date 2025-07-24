package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertDeliveryUnitsWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertDeliveryUnitsWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertDeliveryUnitsWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertDeliveryUnitsWorkflow, error) {
	return UpsertDeliveryUnitsWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertDeliveryUnitsWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertDeliveryUnitsWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertDeliveryUnitsStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.DeliveryUnitsUpserted(), Src: []string{w.UpsertDeliveryUnitsStarted()}, Dst: w.DeliveryUnitsUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertDeliveryUnitsWorkflow) WorkflowName() string {
	return "upsert_delivery_units_workflow"
}

func (w UpsertDeliveryUnitsWorkflow) DeliveryUnitsUpserted() string {
	return "delivery_units_upserted"
}

func (w UpsertDeliveryUnitsWorkflow) UpsertDeliveryUnitsStarted() string {
	return "upsert_delivery_units_started"
}

func (w UpsertDeliveryUnitsWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertDeliveryUnitsWorkflow) CanTransitionToDeliveryUnitsUpserted() bool {
	return w.fsm.Can(w.DeliveryUnitsUpserted())
}

func (w UpsertDeliveryUnitsWorkflow) IsDeliveryUnitsUpserted() bool {
	return w.fsm.Current() == w.DeliveryUnitsUpserted()
}

func (w UpsertDeliveryUnitsWorkflow) IsUpsertDeliveryUnitsStarted() bool {
	return w.fsm.Current() == w.UpsertDeliveryUnitsStarted()
}

// SetDeliveryUnitsUpsertedTransition completa el upsert de delivery units
func (w UpsertDeliveryUnitsWorkflow) SetDeliveryUnitsUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.DeliveryUnitsUpserted())
}