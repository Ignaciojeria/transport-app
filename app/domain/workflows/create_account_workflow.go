package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type CreateAccountWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewCreateAccountWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewCreateAccountWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (CreateAccountWorkflow, error) {
	return CreateAccountWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w CreateAccountWorkflow) Restore(ctx context.Context, idempotencyKey string) (CreateAccountWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.RegistrationStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.AccountCreated(), Src: []string{w.RegistrationStarted()}, Dst: w.AccountCreated()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w CreateAccountWorkflow) WorkflowName() string {
	return "create_account_workflow"
}

func (w CreateAccountWorkflow) AccountCreated() string {
	return "account_created"
}

func (w CreateAccountWorkflow) RegistrationStarted() string {
	return "registration_started"
}

func (w CreateAccountWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w CreateAccountWorkflow) TransitionToAccountCreated(ctx context.Context) error {
	return w.fsm.Event(ctx, w.AccountCreated())
}

func (w CreateAccountWorkflow) CanTransitionToAccountCreated() bool {
	return w.fsm.Can(w.AccountCreated())
}

func (w CreateAccountWorkflow) IsAccountCreated() bool {
	return w.fsm.Current() == w.AccountCreated()
}

func (w CreateAccountWorkflow) IsRegistrationStarted() bool {
	return w.fsm.Current() == w.RegistrationStarted()
}

// CompleteAccountCreation completa la creación de la cuenta
func (w CreateAccountWorkflow) SetAccountCreatedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.AccountCreated())
}
