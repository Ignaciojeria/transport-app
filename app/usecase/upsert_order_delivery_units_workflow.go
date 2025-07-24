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

type UpsertOrderDeliveryUnitsWorkflow func(ctx context.Context, order domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertOrderDeliveryUnitsWorkflow,
		workflows.NewUpsertOrderDeliveryUnitsWorkflow,
		tidbrepository.NewUpsertOrderDeliveryUnits,
		observability.NewObservability)
}

func NewUpsertOrderDeliveryUnitsWorkflow(
	domainWorkflow workflows.UpsertOrderDeliveryUnitsWorkflow,
	upsertOrderDeliveryUnits tidbrepository.UpsertOrderDeliveryUnits,
	obs observability.Observability,
) UpsertOrderDeliveryUnitsWorkflow {
	return func(ctx context.Context, order domain.Order) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "order_delivery_units", order)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetOrderDeliveryUnitsUpsertedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"order_delivery_units_doc_id", order.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertOrderDeliveryUnits(ctx, order, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert order delivery units: %w", err)
		}
		return nil
	}
}