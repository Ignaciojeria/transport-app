package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/restyclient/webhook"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type PublishWebhookWorkflow func(ctx context.Context, body interface{}, webhookType string) error

func init() {
	ioc.Registry(
		NewPublishWebhookWorkflow,
		workflows.NewPublishWebhookWorkflow,
		webhook.NewPostWebhook,
		observability.NewObservability,
	)
}

func NewPublishWebhookWorkflow(
	workflow workflows.PublishWebhookWorkflow,
	postWebhook webhook.PostWebhook,
	obs observability.Observability,
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

		// Intentar publicar el webhook
		err = postWebhook(ctx, body, webhookType)
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
