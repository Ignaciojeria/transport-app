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

type UpsertOrderTypeWorkflow func(ctx context.Context, ot domain.OrderType) error

func init() {
	ioc.Registry(
		NewUpsertOrderTypeWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertOrderType,
		observability.NewObservability)
}

func NewUpsertOrderTypeWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertOrderType tidbrepository.UpsertOrderType,
	obs observability.Observability,
) UpsertOrderTypeWorkflow {
	return func(ctx context.Context, ot domain.OrderType) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		config := workflows.CreateUpsertWorkflow("order_type")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"order_type_doc_id", ot.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertOrderType(ctx, ot, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert order type: %w", err)
		}
		return nil
	}
}