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

type UpsertDeliveryUnitsWorkflow func(ctx context.Context, deliveryUnits []domain.DeliveryUnit) error

func init() {
	ioc.Registry(
		NewUpsertDeliveryUnitsWorkflow,
		workflows.NewUpsertDeliveryUnitsWorkflow,
		tidbrepository.NewUpsertDeliveryUnits,
		observability.NewObservability)
}

func NewUpsertDeliveryUnitsWorkflow(
	domainWorkflow workflows.UpsertDeliveryUnitsWorkflow,
	upsertDeliveryUnits tidbrepository.UpsertDeliveryUnits,
	obs observability.Observability,
) UpsertDeliveryUnitsWorkflow {
	return func(ctx context.Context, deliveryUnits []domain.DeliveryUnit) error {
		// Usar el hash de los delivery units como idempotency key
		key, err := canonicaljson.HashKey(ctx, "delivery_units", deliveryUnits)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetDeliveryUnitsUpsertedTransition(ctx); err != nil {
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
