package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertDeliveryUnitsLabelsWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertDeliveryUnitsLabelsWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertDeliveryUnitsLabelsWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertDeliveryUnitsLabelsWorkflow, error) {
	return UpsertDeliveryUnitsLabelsWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertDeliveryUnitsLabelsWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertDeliveryUnitsLabelsWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertDeliveryUnitsLabelsStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.DeliveryUnitsLabelsUpserted(), Src: []string{w.UpsertDeliveryUnitsLabelsStarted()}, Dst: w.DeliveryUnitsLabelsUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertDeliveryUnitsLabelsWorkflow) WorkflowName() string {
	return "upsert_delivery_units_labels_workflow"
}

func (w UpsertDeliveryUnitsLabelsWorkflow) DeliveryUnitsLabelsUpserted() string {
	return "delivery_units_labels_upserted"
}

func (w UpsertDeliveryUnitsLabelsWorkflow) UpsertDeliveryUnitsLabelsStarted() string {
	return "upsert_delivery_units_labels_started"
}

func (w UpsertDeliveryUnitsLabelsWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertDeliveryUnitsLabelsWorkflow) TransitionToDeliveryUnitsLabelsUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.DeliveryUnitsLabelsUpserted())
}

func (w UpsertDeliveryUnitsLabelsWorkflow) CanTransitionToDeliveryUnitsLabelsUpserted() bool {
	return w.fsm.Can(w.DeliveryUnitsLabelsUpserted())
}

func (w UpsertDeliveryUnitsLabelsWorkflow) IsDeliveryUnitsLabelsUpserted() bool {
	return w.fsm.Current() == w.DeliveryUnitsLabelsUpserted()
}

func (w UpsertDeliveryUnitsLabelsWorkflow) IsUpsertDeliveryUnitsLabelsStarted() bool {
	return w.fsm.Current() == w.UpsertDeliveryUnitsLabelsStarted()
}

// SetDeliveryUnitsLabelsUpsertedTransition completa el upsert de delivery units labels
func (w UpsertDeliveryUnitsLabelsWorkflow) SetDeliveryUnitsLabelsUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.DeliveryUnitsLabelsUpserted())
}