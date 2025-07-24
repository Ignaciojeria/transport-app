package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertDeliveryUnitsHistoryWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertDeliveryUnitsHistoryWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertDeliveryUnitsHistoryWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertDeliveryUnitsHistoryWorkflow, error) {
	return UpsertDeliveryUnitsHistoryWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertDeliveryUnitsHistoryWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertDeliveryUnitsHistoryWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertDeliveryUnitsHistoryStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.DeliveryUnitsHistoryUpserted(), Src: []string{w.UpsertDeliveryUnitsHistoryStarted()}, Dst: w.DeliveryUnitsHistoryUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertDeliveryUnitsHistoryWorkflow) WorkflowName() string {
	return "upsert_delivery_units_history_workflow"
}

func (w UpsertDeliveryUnitsHistoryWorkflow) DeliveryUnitsHistoryUpserted() string {
	return "delivery_units_history_upserted"
}

func (w UpsertDeliveryUnitsHistoryWorkflow) UpsertDeliveryUnitsHistoryStarted() string {
	return "upsert_delivery_units_history_started"
}

func (w UpsertDeliveryUnitsHistoryWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertDeliveryUnitsHistoryWorkflow) CanTransitionToDeliveryUnitsHistoryUpserted() bool {
	return w.fsm.Can(w.DeliveryUnitsHistoryUpserted())
}

func (w UpsertDeliveryUnitsHistoryWorkflow) IsDeliveryUnitsHistoryUpserted() bool {
	return w.fsm.Current() == w.DeliveryUnitsHistoryUpserted()
}

func (w UpsertDeliveryUnitsHistoryWorkflow) IsUpsertDeliveryUnitsHistoryStarted() bool {
	return w.fsm.Current() == w.UpsertDeliveryUnitsHistoryStarted()
}

// SetDeliveryUnitsHistoryUpsertedTransition completa el upsert de delivery units history
func (w UpsertDeliveryUnitsHistoryWorkflow) SetDeliveryUnitsHistoryUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.DeliveryUnitsHistoryUpserted())
}