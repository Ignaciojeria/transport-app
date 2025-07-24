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

type UpsertDeliveryUnitsLabelsWorkflow func(ctx context.Context, order domain.Order) error

func init() {
	ioc.Registry(
		NewUpsertDeliveryUnitsLabelsWorkflow,
		workflows.NewUpsertDeliveryUnitsLabelsWorkflow,
		tidbrepository.NewUpsertDeliveryUnitsLabels,
		observability.NewObservability)
}

func NewUpsertDeliveryUnitsLabelsWorkflow(
	domainWorkflow workflows.UpsertDeliveryUnitsLabelsWorkflow,
	upsertDeliveryUnitsLabels tidbrepository.UpsertDeliveryUnitsLabels,
	obs observability.Observability,
) UpsertDeliveryUnitsLabelsWorkflow {
	return func(ctx context.Context, order domain.Order) error {
		key, err := canonicaljson.HashKey(ctx, "delivery_units_labels", order)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetDeliveryUnitsLabelsUpsertedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"delivery_units_labels_doc_id", order.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertDeliveryUnitsLabels(ctx, order, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert delivery units labels: %w", err)
		}
		return nil
	}
}
