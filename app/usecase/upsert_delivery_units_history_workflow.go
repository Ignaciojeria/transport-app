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

type UpsertDeliveryUnitsHistoryWorkflow func(ctx context.Context, plan domain.Plan) error

func init() {
	ioc.Registry(
		NewUpsertDeliveryUnitsHistoryWorkflow,
		workflows.NewUpsertDeliveryUnitsHistoryWorkflow,
		tidbrepository.NewUpsertDeliveryUnitsHistory,
		observability.NewObservability)
}

func NewUpsertDeliveryUnitsHistoryWorkflow(
	domainWorkflow workflows.UpsertDeliveryUnitsHistoryWorkflow,
	upsertDeliveryUnitsHistory tidbrepository.UpsertDeliveryUnitsHistory,
	obs observability.Observability,
) UpsertDeliveryUnitsHistoryWorkflow {
	return func(ctx context.Context, plan domain.Plan) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetDeliveryUnitsHistoryUpsertedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"delivery_units_history_doc_id", key)
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertDeliveryUnitsHistory(ctx, plan, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert delivery units history: %w", err)
		}
		return nil
	}
}
