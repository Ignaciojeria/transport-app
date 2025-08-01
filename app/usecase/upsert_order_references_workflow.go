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

type UpsertOrderReferencesWorkflow func(ctx context.Context, order domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertOrderReferencesWorkflow,
		workflows.NewUpsertOrderReferencesWorkflow,
		tidbrepository.NewUpsertOrderReferences,
		observability.NewObservability)
}

func NewUpsertOrderReferencesWorkflow(
	domainWorkflow workflows.UpsertOrderReferencesWorkflow,
	upsertOrderReferences tidbrepository.UpsertOrderReferences,
	obs observability.Observability,
) UpsertOrderReferencesWorkflow {
	return func(ctx context.Context, order domain.Order) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetOrderReferencesUpsertedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"order_references_doc_id", order.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertOrderReferences(ctx, order, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert order references: %w", err)
		}
		return nil
	}
}
