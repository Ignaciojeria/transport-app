package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertClientCredentialsWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertClientCredentialsWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertClientCredentialsWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertClientCredentialsWorkflow, error) {
	return UpsertClientCredentialsWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertClientCredentialsWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertClientCredentialsWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertClientCredentialsStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.ClientCredentialsUpserted(), Src: []string{w.UpsertClientCredentialsStarted()}, Dst: w.ClientCredentialsUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertClientCredentialsWorkflow) WorkflowName() string {
	return "upsert_client_credentials_workflow"
}

func (w UpsertClientCredentialsWorkflow) ClientCredentialsUpserted() string {
	return "client_credentials_upserted"
}

func (w UpsertClientCredentialsWorkflow) UpsertClientCredentialsStarted() string {
	return "upsert_client_credentials_started"
}

func (w UpsertClientCredentialsWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertClientCredentialsWorkflow) CanTransitionToClientCredentialsUpserted() bool {
	return w.fsm.Can(w.ClientCredentialsUpserted())
}

func (w UpsertClientCredentialsWorkflow) IsClientCredentialsUpserted() bool {
	return w.fsm.Current() == w.ClientCredentialsUpserted()
}

func (w UpsertClientCredentialsWorkflow) IsUpsertClientCredentialsStarted() bool {
	return w.fsm.Current() == w.UpsertClientCredentialsStarted()
}

// SetClientCredentialsUpsertedTransition completa el upsert de client credentials
func (w UpsertClientCredentialsWorkflow) SetClientCredentialsUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.ClientCredentialsUpserted())
}