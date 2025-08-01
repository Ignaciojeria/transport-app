package workflows

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type PublishWebhookWorkflow struct {
	IdempotencyKey       string
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey
	fsm                  *fsm.FSM
}

func init() {
	ioc.Registry(NewPublishWebhookWorkflow,
		tidbrepository.NewGetLastFSMTransitionByIdempotencyKey)
}

func NewPublishWebhookWorkflow(
	getLastFSMTransition tidbrepository.GetLastFSMTransitionByIdempotencyKey) (PublishWebhookWorkflow, error) {
	return PublishWebhookWorkflow{
		getLastFSMTransition: getLastFSMTransition,
	}, nil
}

func (w PublishWebhookWorkflow) Restore(ctx context.Context, idempotencyKey string) (PublishWebhookWorkflow, error) {
	lastTransition, err := w.getLastFSMTransition(ctx, idempotencyKey, w.WorkflowName())
	w.IdempotencyKey = idempotencyKey
	if err != nil {
		return w, err
	}
	transition := lastTransition.State
	if transition == "" {
		transition = w.PublishWebhookStarted()
	}
	w.fsm = fsm.NewFSM(
		transition,
		fsm.Events{
			{Name: w.WebhookPublished(), Src: []string{w.PublishWebhookStarted()}, Dst: w.WebhookPublished()},
			{Name: w.WebhookFailed(), Src: []string{w.PublishWebhookStarted()}, Dst: w.WebhookFailed()},
		},
		fsm.Callbacks{},
	)
	return w, nil
}

func (w PublishWebhookWorkflow) WorkflowName() string {
	return "publish_webhook_workflow"
}

func (w PublishWebhookWorkflow) WebhookPublished() string {
	return "webhook_published"
}

func (w PublishWebhookWorkflow) WebhookFailed() string {
	return "webhook_failed"
}

func (w PublishWebhookWorkflow) PublishWebhookStarted() string {
	return "publish_webhook_started"
}

func (w PublishWebhookWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		State:          w.fsm.Current(),
	}
}

func (w PublishWebhookWorkflow) CanTransitionToWebhookPublished() bool {
	return w.fsm.Can(w.WebhookPublished())
}

func (w PublishWebhookWorkflow) CanTransitionToWebhookFailed() bool {
	return w.fsm.Can(w.WebhookFailed())
}

func (w PublishWebhookWorkflow) IsWebhookPublished() bool {
	return w.fsm.Current() == w.WebhookPublished()
}

func (w PublishWebhookWorkflow) IsWebhookFailed() bool {
	return w.fsm.Current() == w.WebhookFailed()
}

func (w PublishWebhookWorkflow) IsPublishWebhookStarted() bool {
	return w.fsm.Current() == w.PublishWebhookStarted()
}

// SetWebhookPublishedTransition marca el webhook como publicado exitosamente
func (w PublishWebhookWorkflow) SetWebhookPublishedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.WebhookPublished())
}

// SetWebhookFailedTransition marca el webhook como fallido
func (w PublishWebhookWorkflow) SetWebhookFailedTransition(ctx context.Context) error {
	return w.fsm.Event(ctx, w.WebhookFailed())
}
