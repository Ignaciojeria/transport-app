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

type UpsertOrderHeadersWorkflow func(ctx context.Context, headers domain.Headers) error

func init() {
	ioc.Registry(
		NewUpsertOrderHeadersWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertOrderHeaders,
		observability.NewObservability)
}

func NewUpsertOrderHeadersWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertOrderHeaders tidbrepository.UpsertOrderHeaders,
	obs observability.Observability,
) UpsertOrderHeadersWorkflow {
	return func(ctx context.Context, headers domain.Headers) error {
		// Usar el documentID como idempotency key para el workflow
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		config := workflows.CreateUpsertWorkflow("order_headers")
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"headers_doc_id", headers.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertOrderHeaders(ctx, headers, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert order headers: %w", err)
		}
		return nil
	}
}
