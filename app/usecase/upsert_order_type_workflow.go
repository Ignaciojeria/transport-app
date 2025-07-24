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

type UpsertOrderTypeWorkflow func(ctx context.Context, ot domain.OrderType) error

func init() {
	ioc.Registry(
		NewUpsertOrderTypeWorkflow,
		workflows.NewUpsertOrderTypeWorkflow,
		tidbrepository.NewUpsertOrderType,
		observability.NewObservability)
}

func NewUpsertOrderTypeWorkflow(
	domainWorkflow workflows.UpsertOrderTypeWorkflow,
	upsertOrderType tidbrepository.UpsertOrderType,
	obs observability.Observability,
) UpsertOrderTypeWorkflow {
	return func(ctx context.Context, ot domain.OrderType) error {
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "order_type", ot)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetOrderTypeUpsertedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"order_type_doc_id", ot.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertOrderType(ctx, ot, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert order type: %w", err)
		}
		return nil
	}
}