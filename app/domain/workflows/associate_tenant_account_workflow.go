package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type AssociateTenantAccountWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewAssociateTenantAccountWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewAssociateTenantAccountWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (AssociateTenantAccountWorkflow, error) {
	return AssociateTenantAccountWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w AssociateTenantAccountWorkflow) Restore(ctx context.Context, idempotencyKey string) (AssociateTenantAccountWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.AssociationStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.AssociationCompleted(), Src: []string{w.AssociationStarted()}, Dst: w.AssociationCompleted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w AssociateTenantAccountWorkflow) WorkflowName() string {
	return "associate_tenant_account_workflow"
}

func (w AssociateTenantAccountWorkflow) AssociationCompleted() string {
	return "association_completed"
}

func (w AssociateTenantAccountWorkflow) AssociationStarted() string {
	return "association_started"
}

func (w AssociateTenantAccountWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w AssociateTenantAccountWorkflow) TransitionToAssociationCompleted(ctx context.Context) error {
	return w.fsm.Event(ctx, w.AssociationCompleted())
}

func (w AssociateTenantAccountWorkflow) CanTransitionToAssociationCompleted() bool {
	return w.fsm.Can(w.AssociationCompleted())
}

func (w AssociateTenantAccountWorkflow) IsAssociationCompleted() bool {
	return w.fsm.Current() == w.AssociationCompleted()
}

func (w AssociateTenantAccountWorkflow) IsAssociationStarted() bool {
	return w.fsm.Current() == w.AssociationStarted()
}

// CompleteAssociation completa la asociación del tenant con el account
func (w AssociateTenantAccountWorkflow) SetAssociationCompletedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.AssociationCompleted())
}
