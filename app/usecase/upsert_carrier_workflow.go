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

type UpsertCarrierWorkflow func(ctx context.Context, c domain.Carrier) error

func init() {
	ioc.Registry(
		NewUpsertCarrierWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertCarrier,
		observability.NewObservability)
}

func NewUpsertCarrierWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertCarrier tidbrepository.UpsertCarrier,
	obs observability.Observability,
) UpsertCarrierWorkflow {
	return func(ctx context.Context, c domain.Carrier) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "carrier", c)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		config := workflows.CreateUpsertWorkflow("carrier")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"carrier_doc_id", c.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertCarrier(ctx, c, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert carrier: %w", err)
		}
		return nil
	}
}