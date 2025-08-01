package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertWebhookWorkflow func(ctx context.Context, w domain.Webhook) error

func init() {
	ioc.Registry(
		NewUpsertWebhookWorkflow,
		workflows.NewUpsertWebhookWorkflow,
		tidbrepository.NewUpsertWebhook,
		observability.NewObservability)
}

func NewUpsertWebhookWorkflow(
	domainWorkflow workflows.UpsertWebhookWorkflow,
	upsertWebhook tidbrepository.UpsertWebhook,
	obs observability.Observability,
) UpsertWebhookWorkflow {
	return func(ctx context.Context, w domain.Webhook) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}

		if err := workflow.SetWebhookUpsertedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"webhook_doc_id", w.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)

		err = upsertWebhook(ctx, w, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert webhook: %w", err)
		}
		return nil
	}
}
