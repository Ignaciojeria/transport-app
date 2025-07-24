package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertCarrierWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertCarrierWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertCarrierWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertCarrierWorkflow, error) {
	return UpsertCarrierWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertCarrierWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertCarrierWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertCarrierStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.CarrierUpserted(), Src: []string{w.UpsertCarrierStarted()}, Dst: w.CarrierUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertCarrierWorkflow) WorkflowName() string {
	return "upsert_carrier_workflow"
}

func (w UpsertCarrierWorkflow) CarrierUpserted() string {
	return "carrier_upserted"
}

func (w UpsertCarrierWorkflow) UpsertCarrierStarted() string {
	return "upsert_carrier_started"
}

func (w UpsertCarrierWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertCarrierWorkflow) CanTransitionToCarrierUpserted() bool {
	return w.fsm.Can(w.CarrierUpserted())
}

func (w UpsertCarrierWorkflow) IsCarrierUpserted() bool {
	return w.fsm.Current() == w.CarrierUpserted()
}

func (w UpsertCarrierWorkflow) IsUpsertCarrierStarted() bool {
	return w.fsm.Current() == w.UpsertCarrierStarted()
}

// SetCarrierUpsertedTransition completa el upsert de carrier
func (w UpsertCarrierWorkflow) SetCarrierUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.CarrierUpserted())
}