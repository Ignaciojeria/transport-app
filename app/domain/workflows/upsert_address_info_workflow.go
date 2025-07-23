package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertAddressInfoWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertAddressInfoWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertAddressInfoWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertAddressInfoWorkflow, error) {
	return UpsertAddressInfoWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertAddressInfoWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertAddressInfoWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertAddressInfoStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.AddressInfoUpserted(), Src: []string{w.UpsertAddressInfoStarted()}, Dst: w.AddressInfoUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertAddressInfoWorkflow) WorkflowName() string {
	return "upsert_address_info_workflow"
}

func (w UpsertAddressInfoWorkflow) AddressInfoUpserted() string {
	return "address_info_upserted"
}

func (w UpsertAddressInfoWorkflow) UpsertAddressInfoStarted() string {
	return "upsert_address_info_started"
}

func (w UpsertAddressInfoWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertAddressInfoWorkflow) TransitionToAddressInfoUpserted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.AddressInfoUpserted())
}

func (w UpsertAddressInfoWorkflow) CanTransitionToAddressInfoUpserted() bool {
	return w.fsm.Can(w.AddressInfoUpserted())
}

func (w UpsertAddressInfoWorkflow) IsAddressInfoUpserted() bool {
	return w.fsm.Current() == w.AddressInfoUpserted()
}

func (w UpsertAddressInfoWorkflow) IsUpsertAddressInfoStarted() bool {
	return w.fsm.Current() == w.UpsertAddressInfoStarted()
}

// SetAddressInfoUpsertedTransition completa el upsert de address info
func (w UpsertAddressInfoWorkflow) SetAddressInfoUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.AddressInfoUpserted())
}
