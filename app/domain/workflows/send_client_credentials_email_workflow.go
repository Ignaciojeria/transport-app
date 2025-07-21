package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type SendClientCredentialsEmailWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewSendClientCredentialsEmailWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewSendClientCredentialsEmailWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (SendClientCredentialsEmailWorkflow, error) {
	return SendClientCredentialsEmailWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w SendClientCredentialsEmailWorkflow) Restore(ctx context.Context, idempotencyKey string) (SendClientCredentialsEmailWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.EmailSendingStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.EmailSent(), Src: []string{w.EmailSendingStarted()}, Dst: w.EmailSent()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w SendClientCredentialsEmailWorkflow) WorkflowName() string {
	return "send_client_credentials_email_workflow"
}

func (w SendClientCredentialsEmailWorkflow) EmailSent() string {
	return "email_sent"
}

func (w SendClientCredentialsEmailWorkflow) EmailSendingStarted() string {
	return "email_sending_started"
}

func (w SendClientCredentialsEmailWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w SendClientCredentialsEmailWorkflow) TransitionToEmailSent(ctx context.Context) error {
	return w.fsm.Event(ctx, w.EmailSent())
}

func (w SendClientCredentialsEmailWorkflow) CanTransitionToEmailSent() bool {
	return w.fsm.Can(w.EmailSent())
}

func (w SendClientCredentialsEmailWorkflow) IsEmailSent() bool {
	return w.fsm.Current() == w.EmailSent()
}

func (w SendClientCredentialsEmailWorkflow) IsEmailSendingStarted() bool {
	return w.fsm.Current() == w.EmailSendingStarted()
}

// CompleteEmailSending completa el envío del email de credenciales
func (w SendClientCredentialsEmailWorkflow) SetEmailSentTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.EmailSent())
}
