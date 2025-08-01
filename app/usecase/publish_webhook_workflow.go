package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/restyclient/webhook"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type PublishWebhookWorkflow func(ctx context.Context, webhookType string) error

func init() {
	ioc.Registry(
		NewPublishWebhookWorkflow,
		workflows.NewPublishWebhookWorkflow,
		webhook.NewPostWebhook,
		tidbrepository.NewSaveFSMTransition,
		observability.NewObservability,
		tidbrepository.NewFindWebhookByDocumentID)
}

func NewPublishWebhookWorkflow(
	workflow workflows.PublishWebhookWorkflow,
	postWebhook webhook.PostWebhook,
	saveFSMTransition tidbrepository.SaveFSMTransition,
	obs observability.Observability,
	findWebhookByDocumentID tidbrepository.FindWebhookByDocumentID,
) PublishWebhookWorkflow {
	return func(ctx context.Context, webhookType string) error {
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
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"webhook_type", webhookType)
			return nil
		}

		webhook, err := findWebhookByDocumentID(ctx, domain.Webhook{Type: webhookType}.DocID(ctx))
		if err != nil {
			return fmt.Errorf("failed to find webhook by document ID: %w", err)
		}
		// Intentar publicar el webhook
		err = postWebhook(ctx, webhook)
		if err != nil {
			return fmt.Errorf("failed to publish webhook: %w", err)
		}

		fsmState := workflowInstance.Map(ctx)
		err = saveFSMTransition(ctx, fsmState)
		if err != nil {
			return fmt.Errorf("failed to save FSM transition: %w", err)
		}

		return nil
	}
}
