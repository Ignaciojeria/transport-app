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
		workflows.NewUpsertVehicleWorkflow,
		tidbrepository.NewUpsertVehicle,
		observability.NewObservability)
}

func NewUpsertVehicleWorkflow(
	domainWorkflow workflows.UpsertVehicleWorkflow,
	upsertVehicle tidbrepository.UpsertVehicle,
	obs observability.Observability,
) UpsertVehicleWorkflow {
	return func(ctx context.Context, v domain.Vehicle) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "vehicle", v)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetVehicleUpsertedTransition(ctx); err != nil {
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