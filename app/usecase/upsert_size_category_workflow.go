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

type UpsertSizeCategoryWorkflow func(ctx context.Context, sc domain.SizeCategory) error

func init() {
	ioc.Registry(
		NewUpsertSizeCategoryWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertSizeCategory,
		observability.NewObservability)
}

func NewUpsertSizeCategoryWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertSizeCategory tidbrepository.UpsertSizeCategory,
	obs observability.Observability,
) UpsertSizeCategoryWorkflow {
	return func(ctx context.Context, sc domain.SizeCategory) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		config := workflows.CreateUpsertWorkflow("size_category")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"size_category_doc_id", sc.DocumentID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertSizeCategory(ctx, sc, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert size category: %w", err)
		}
		return nil
	}
}