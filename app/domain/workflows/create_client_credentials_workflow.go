package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type CreateClientCredentialsWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewCreateClientCredentialsWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewCreateClientCredentialsWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (CreateClientCredentialsWorkflow, error) {
	return CreateClientCredentialsWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w CreateClientCredentialsWorkflow) Restore(ctx context.Context, idempotencyKey string) (CreateClientCredentialsWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.CredentialsCreationStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.CredentialsCreated(), Src: []string{w.CredentialsCreationStarted()}, Dst: w.CredentialsCreated()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w CreateClientCredentialsWorkflow) WorkflowName() string {
	return "create_client_credentials_workflow"
}

func (w CreateClientCredentialsWorkflow) CredentialsCreated() string {
	return "credentials_created"
}

func (w CreateClientCredentialsWorkflow) CredentialsCreationStarted() string {
	return "credentials_creation_started"
}

func (w CreateClientCredentialsWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w CreateClientCredentialsWorkflow) TransitionToCredentialsCreated(ctx context.Context) error {
	return w.fsm.Event(ctx, w.CredentialsCreated())
}

func (w CreateClientCredentialsWorkflow) CanTransitionToCredentialsCreated() bool {
	return w.fsm.Can(w.CredentialsCreated())
}

func (w CreateClientCredentialsWorkflow) IsCredentialsCreated() bool {
	return w.fsm.Current() == w.CredentialsCreated()
}

func (w CreateClientCredentialsWorkflow) IsCredentialsCreationStarted() bool {
	return w.fsm.Current() == w.CredentialsCreationStarted()
}

// CompleteCredentialsCreation completa la creación de las credenciales
func (w CreateClientCredentialsWorkflow) SetCredentialsCreatedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.CredentialsCreated())
}
