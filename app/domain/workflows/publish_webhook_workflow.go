package workflows

import (
	"context"
	"encoding/json"
	"errors"
	"transport-app/app/adapter/out/storjbucket"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/looplab/fsm"
	"go.opentelemetry.io/otel/trace"
)

type PublishWebhookWorkflow struct {
	IdempotencyKey string
	NextInput      []byte
	storjBucket    *storjbucket.TransportAppBucket
	fsm            *fsm.FSM
}

func init() {
	ioc.Registry(NewPublishWebhookWorkflow,
		storjbucket.NewTransportAppBucket)
}

func NewPublishWebhookWorkflow(
	storjBucket *storjbucket.TransportAppBucket,
) (PublishWebhookWorkflow, error) {
	return PublishWebhookWorkflow{
		storjBucket: storjBucket,
	}, nil
}

func (w PublishWebhookWorkflow) Restore(ctx context.Context, idempotencyKey string) (PublishWebhookWorkflow, error) {
	w.IdempotencyKey = idempotencyKey

	// Obtener token desde el contexto
	token, ok := sharedcontext.BucketTokenFromContext(ctx)
	if !ok {
		return w, errors.New("token del bucket no encontrado en el contexto")
	}

	// Intentar recuperar el estado desde StorJ bucket
	data, err := w.storjBucket.DownloadWithToken(ctx, token, idempotencyKey)

	var transition string
	var nextInput []byte

	if err != nil {
		// Si no existe el estado, usar el estado inicial
		transition = w.PublishWebhookStarted()
		nextInput = nil
	} else {
		// Deserializar el estado guardado
		var fsmState domain.FSMState
		if err := json.Unmarshal(data, &fsmState); err != nil {
			return w, err
		}
		transition = fsmState.State
		nextInput = fsmState.NextInput
	}

	// Si no hay estado guardado, usar el estado inicial
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
	w.NextInput = nextInput
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

func (w PublishWebhookWorkflow) SaveState(ctx context.Context) error {
	// Obtener token desde el contexto
	token, ok := sharedcontext.BucketTokenFromContext(ctx)
	if !ok {
		return errors.New("token del bucket no encontrado en el contexto")
	}

	// Serializar el estado actual
	fsmState := w.Map(ctx)
	data, err := json.Marshal(fsmState)
	if err != nil {
		return err
	}

	// Guardar en StorJ bucket usando el mismo patrón que en el publisher
	return w.storjBucket.UploadWithToken(ctx, token, w.IdempotencyKey, data)
}

func (w PublishWebhookWorkflow) Map(ctx context.Context) domain.FSMState {
	return domain.FSMState{
		Workflow:       w.WorkflowName(),
		TraceID:        trace.SpanContextFromContext(ctx).TraceID().String(),
		IdempotencyKey: w.IdempotencyKey,
		NextInput:      w.NextInput,
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
