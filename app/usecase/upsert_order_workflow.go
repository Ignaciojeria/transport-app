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

type UpsertOrderWorkflow func(ctx context.Context, o domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertOrderWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertOrder,
		observability.NewObservability)
}

func NewUpsertOrderWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertOrder tidbrepository.UpsertOrder,
	obs observability.Observability,
) UpsertOrderWorkflow {
	return func(ctx context.Context, o domain.Order) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		config := workflows.CreateUpsertWorkflow("order")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"order_doc_id", o.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertOrder(ctx, o, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert order: %w", err)
		}
		return nil
	}
}