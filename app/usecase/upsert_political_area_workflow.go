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

type UpsertPoliticalAreaWorkflow func(ctx context.Context, pa domain.PoliticalArea) error

func init() {
	ioc.Registry(
		NewUpsertPoliticalAreaWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertPoliticalArea,
		observability.NewObservability)
}

func NewUpsertPoliticalAreaWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertPoliticalArea tidbrepository.UpsertPoliticalArea,
	obs observability.Observability,
) UpsertPoliticalAreaWorkflow {
	return func(ctx context.Context, pa domain.PoliticalArea) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "political_area", pa)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		config := workflows.CreateUpsertWorkflow("political_area")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"political_area_doc_id", pa.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertPoliticalArea(ctx, pa, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert political area: %w", err)
		}
		return nil
	}
}