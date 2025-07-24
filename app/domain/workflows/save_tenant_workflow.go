package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type SaveTenantWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewSaveTenantWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewSaveTenantWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (SaveTenantWorkflow, error) {
	return SaveTenantWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w SaveTenantWorkflow) Restore(ctx context.Context, idempotencyKey string) (SaveTenantWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.SaveTenantStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.TenantSaved(), Src: []string{w.SaveTenantStarted()}, Dst: w.TenantSaved()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w SaveTenantWorkflow) WorkflowName() string {
	return "save_tenant_workflow"
}

func (w SaveTenantWorkflow) TenantSaved() string {
	return "tenant_saved"
}

func (w SaveTenantWorkflow) SaveTenantStarted() string {
	return "save_tenant_started"
}

func (w SaveTenantWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w SaveTenantWorkflow) CanTransitionToTenantSaved() bool {
	return w.fsm.Can(w.TenantSaved())
}

func (w SaveTenantWorkflow) IsTenantSaved() bool {
	return w.fsm.Current() == w.TenantSaved()
}

func (w SaveTenantWorkflow) IsSaveTenantStarted() bool {
	return w.fsm.Current() == w.SaveTenantStarted()
}

// SetTenantSavedTransition completa el save de tenant
func (w SaveTenantWorkflow) SetTenantSavedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.TenantSaved())
}