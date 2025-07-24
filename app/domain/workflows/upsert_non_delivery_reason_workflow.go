package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertNonDeliveryReasonWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertNonDeliveryReasonWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertNonDeliveryReasonWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertNonDeliveryReasonWorkflow, error) {
	return UpsertNonDeliveryReasonWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertNonDeliveryReasonWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertNonDeliveryReasonWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertNonDeliveryReasonStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.NonDeliveryReasonUpserted(), Src: []string{w.UpsertNonDeliveryReasonStarted()}, Dst: w.NonDeliveryReasonUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertNonDeliveryReasonWorkflow) WorkflowName() string {
	return "upsert_non_delivery_reason_workflow"
}

func (w UpsertNonDeliveryReasonWorkflow) NonDeliveryReasonUpserted() string {
	return "non_delivery_reason_upserted"
}

func (w UpsertNonDeliveryReasonWorkflow) UpsertNonDeliveryReasonStarted() string {
	return "upsert_non_delivery_reason_started"
}

func (w UpsertNonDeliveryReasonWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertNonDeliveryReasonWorkflow) TransitionToNonDeliveryReasonUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.NonDeliveryReasonUpserted())
}

func (w UpsertNonDeliveryReasonWorkflow) CanTransitionToNonDeliveryReasonUpserted() bool {
	return w.fsm.Can(w.NonDeliveryReasonUpserted())
}

func (w UpsertNonDeliveryReasonWorkflow) IsNonDeliveryReasonUpserted() bool {
	return w.fsm.Current() == w.NonDeliveryReasonUpserted()
}

func (w UpsertNonDeliveryReasonWorkflow) IsUpsertNonDeliveryReasonStarted() bool {
	return w.fsm.Current() == w.UpsertNonDeliveryReasonStarted()
}

// SetNonDeliveryReasonUpsertedTransition completa el upsert de non delivery reason
func (w UpsertNonDeliveryReasonWorkflow) SetNonDeliveryReasonUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.NonDeliveryReasonUpserted())
}