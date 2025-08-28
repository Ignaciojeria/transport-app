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

type UpsertVehicleWorkflow func(ctx context.Context, v domain.Vehicle) error

func init() {
	ioc.Registry(
		NewUpsertVehicleWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertVehicle,
		observability.NewObservability)
}

func NewUpsertVehicleWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertVehicle tidbrepository.UpsertVehicle,
	obs observability.Observability,
) UpsertVehicleWorkflow {
	return func(ctx context.Context, v domain.Vehicle) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "vehicle", v)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		
		// Configurar el workflow para vehicle upsert
		config := workflows.CreateUpsertWorkflow("vehicle")
		
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"vehicle_doc_id", v.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertVehicle(ctx, v, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert vehicle: %w", err)
		}
		return nil
	}
}