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

type UpsertAddressInfoWorkflow func(ctx context.Context, ai domain.AddressInfo) error

func init() {
	ioc.Registry(
		NewUpsertAddressInfoWorkflow,
		workflows.NewUpsertAddressInfoWorkflow,
		tidbrepository.NewUpsertAddressInfo,
		observability.NewObservability)
}

func NewUpsertAddressInfoWorkflow(
	domainWorkflow workflows.UpsertAddressInfoWorkflow,
	upsertAddressInfo tidbrepository.UpsertAddressInfo,
	obs observability.Observability,
) UpsertAddressInfoWorkflow {
	return func(ctx context.Context, ai domain.AddressInfo) error {
		// Usar el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		workflow, err := domainWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetAddressInfoUpsertedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"address_doc_id", ai.DocID(ctx).String())
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertAddressInfo(ctx, ai, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert address info: %w", err)
		}
		return nil
	}
}
