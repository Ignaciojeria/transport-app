package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type CreateTenantWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewCreateTenantWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewCreateTenantWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (CreateTenantWorkflow, error) {
	return CreateTenantWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w CreateTenantWorkflow) Restore(ctx context.Context, idempotencyKey string) (CreateTenantWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.TenantCreationStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.TenantCreated(), Src: []string{w.TenantCreationStarted()}, Dst: w.TenantCreated()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w CreateTenantWorkflow) WorkflowName() string {
	return "create_tenant_workflow"
}

func (w CreateTenantWorkflow) TenantCreated() string {
	return "tenant_created"
}

func (w CreateTenantWorkflow) TenantCreationStarted() string {
	return "tenant_creation_started"
}

func (w CreateTenantWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w CreateTenantWorkflow) TransitionToTenantCreated(ctx context.Context) error {
	return w.fsm.Event(ctx, w.TenantCreated())
}

func (w CreateTenantWorkflow) CanTransitionToTenantCreated() bool {
	return w.fsm.Can(w.TenantCreated())
}

func (w CreateTenantWorkflow) IsTenantCreated() bool {
	return w.fsm.Current() == w.TenantCreated()
}

func (w CreateTenantWorkflow) IsTenantCreationStarted() bool {
	return w.fsm.Current() == w.TenantCreationStarted()
}

// CompleteTenantCreation completa la creación del tenant
func (w CreateTenantWorkflow) SetTenantCreatedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.TenantCreated())
}
