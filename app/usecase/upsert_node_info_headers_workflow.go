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

type UpsertNodeInfoHeadersWorkflow func(ctx context.Context, nih domain.Headers) error

func init() {
	ioc.Registry(
		NewUpsertNodeInfoHeadersWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertNodeInfoHeaders,
		observability.NewObservability)
}

func NewUpsertNodeInfoHeadersWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertNodeInfoHeaders tidbrepository.UpsertNodeInfoHeaders,
	obs observability.Observability,
) UpsertNodeInfoHeadersWorkflow {
	return func(ctx context.Context, nih domain.Headers) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "node_info_headers", nih)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		config := workflows.CreateUpsertWorkflow("node_info_headers")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"node_info_headers_doc_id", nih.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertNodeInfoHeaders(ctx, nih, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert node info headers: %w", err)
		}
		return nil
	}
}