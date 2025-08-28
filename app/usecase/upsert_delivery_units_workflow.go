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

type UpsertDeliveryUnitsWorkflow func(ctx context.Context, deliveryUnits []domain.DeliveryUnit) error

func init() {
	ioc.Registry(
		NewUpsertDeliveryUnitsWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertDeliveryUnits,
		observability.NewObservability)
}

func NewUpsertDeliveryUnitsWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertDeliveryUnits tidbrepository.UpsertDeliveryUnits,
	obs observability.Observability,
) UpsertDeliveryUnitsWorkflow {
	return func(ctx context.Context, deliveryUnits []domain.DeliveryUnit) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		config := workflows.CreateUpsertWorkflow("delivery_units")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"delivery_units_doc_id", key)
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertDeliveryUnits(ctx, deliveryUnits, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert delivery units: %w", err)
		}
		return nil
	}
}
