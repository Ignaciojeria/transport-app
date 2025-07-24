package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"
	"transport-app/app/domain/workflows"
	canonicaljson "transport-app/app/shared/caonincaljson"
	"transport-app/app/shared/infrastructure/observability"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertSizeCategoryWorkflow func(ctx context.Context, sc domain.SizeCategory) error

func init() {
	ioc.Registry(
		NewUpsertSizeCategoryWorkflow,
		workflows.NewUpsertSizeCategoryWorkflow,
		tidbrepository.NewUpsertSizeCategory,
		observability.NewObservability)
}

func NewUpsertSizeCategoryWorkflow(
	domainWorkflow workflows.UpsertSizeCategoryWorkflow,
	upsertSizeCategory tidbrepository.UpsertSizeCategory,
	obs observability.Observability,
) UpsertSizeCategoryWorkflow {
	return func(ctx context.Context, sc domain.SizeCategory) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "size_category", sc)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetSizeCategoryUpsertedTransition(ctx); err != nil {
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