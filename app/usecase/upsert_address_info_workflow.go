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
		// Usar el documentID como idempotency key para el workflow
		key, err := canonicaljson.HashKey(ctx, "address_info", ai)
		if err != nil {
			return fmt.Errorf("failed to hash key: %w", err)
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
