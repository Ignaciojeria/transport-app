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

type UpsertVehicleCategoryWorkflow func(ctx context.Context, vc domain.VehicleCategory) error

func init() {
	ioc.Registry(
		NewUpsertVehicleCategoryWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertVehicleCategory,
		observability.NewObservability)
}

func NewUpsertVehicleCategoryWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertVehicleCategory tidbrepository.UpsertVehicleCategory,
	obs observability.Observability,
) UpsertVehicleCategoryWorkflow {
	return func(ctx context.Context, vc domain.VehicleCategory) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "vehicle_category", vc)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		config := workflows.CreateUpsertWorkflow("vehicle_category")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"vehicle_category_doc_id", vc.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertVehicleCategory(ctx, vc, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert vehicle category: %w", err)
		}
		return nil
	}
}