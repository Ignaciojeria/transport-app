package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	client "transport-app/app/adapter/out/restyclient/webhook"
	"transport-app/app/domain"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/infrastructure/natsconn"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/nats-io/nats.go/jetstream"
)

type PublishWebhookWorkflow func(ctx context.Context, body interface{}, webhookType string) error

func init() {
	ioc.Registry(
		NewPublishWebhookWorkflow,
		workflows.NewPublishWebhookWorkflow,
		client.NewPostWebhook,
		observability.NewObservability,
		natsconn.NewKeyValue,
	)
}

func NewPublishWebhookWorkflow(
	workflow workflows.PublishWebhookWorkflow,
	postWebhook client.PostWebhook,
	obs observability.Observability,
	kv jetstream.KeyValue,
) PublishWebhookWorkflow {
	return func(ctx context.Context, body interface{}, webhookType string) error {
		// Obtener el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}

		workflowInstance, err := workflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}

		// Intentar transición a webhook publicado
		if err := workflowInstance.SetWebhookPublishedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx, err.Error())
			return nil
		}

		wh := domain.Webhook{
			Type: webhookType,
		}

		bytes, err := kv.Get(ctx, wh.DocID(ctx).String())
		if err != nil {
			obs.Logger.Error("Error obteniendo webhook", "error", err)
			return err
		}

		var webhook domain.Webhook
		if err := json.Unmarshal(bytes.Value(), &webhook); err != nil {
			obs.Logger.Error("Error deserializando webhook", "error", err)
			return err
		}

		webhook.Body = body

		// Intentar publicar el webhook
		err = postWebhook(ctx, webhook)
		if err != nil {
			return fmt.Errorf("failed to publish webhook: %w", err)
		}

		// Guardar el estado usando el nuevo patrón
		err = workflowInstance.SaveState(ctx)
		if err != nil {
			return fmt.Errorf("failed to save workflow state: %w", err)
		}

		return nil
	}
}
