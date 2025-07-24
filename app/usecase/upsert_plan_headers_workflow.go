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

type UpsertPlanHeadersWorkflow func(ctx context.Context, headers domain.Headers) error

func init() {
	ioc.Registry(
		NewUpsertPlanHeadersWorkflow,
		workflows.NewUpsertPlanHeadersWorkflow,
		tidbrepository.NewUpsertPlanHeaders,
		observability.NewObservability)
}

func NewUpsertPlanHeadersWorkflow(
	domainWorkflow workflows.UpsertPlanHeadersWorkflow,
	upsertPlanHeaders tidbrepository.UpsertPlanHeaders,
	obs observability.Observability,
) UpsertPlanHeadersWorkflow {
	return func(ctx context.Context, headers domain.Headers) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "plan_headers", headers)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetPlanHeadersUpsertedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"plan_headers_doc_id", headers.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertPlanHeaders(ctx, headers, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert plan headers: %w", err)
		}
		return nil
	}
}