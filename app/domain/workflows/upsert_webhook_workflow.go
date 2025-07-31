package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type UpsertWebhookWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewUpsertWebhookWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewUpsertWebhookWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (UpsertWebhookWorkflow, error) {
	return UpsertWebhookWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w UpsertWebhookWorkflow) Restore(ctx context.Context, idempotencyKey string) (UpsertWebhookWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.UpsertWebhookStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.WebhookUpserted(), Src: []string{w.UpsertWebhookStarted()}, Dst: w.WebhookUpserted()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w UpsertWebhookWorkflow) WorkflowName() string {
	return "upsert_webhook_workflow"
}

func (w UpsertWebhookWorkflow) WebhookUpserted() string {
	return "webhook_upserted"
}

func (w UpsertWebhookWorkflow) UpsertWebhookStarted() string {
	return "upsert_webhook_started"
}

func (w UpsertWebhookWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w UpsertWebhookWorkflow) CanTransitionToWebhookUpserted() bool {
	return w.fsm.Can(w.WebhookUpserted())
}

func (w UpsertWebhookWorkflow) IsWebhookUpserted() bool {
	return w.fsm.Current() == w.WebhookUpserted()
}

func (w UpsertWebhookWorkflow) IsUpsertWebhookStarted() bool {
	return w.fsm.Current() == w.UpsertWebhookStarted()
}

// SetWebhookUpsertedTransition completa el upsert de webhook
func (w UpsertWebhookWorkflow) SetWebhookUpsertedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.WebhookUpserted())
}