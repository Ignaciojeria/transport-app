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

type AssociateTenantAccountWorkflow func(ctx context.Context, input domain.TenantAccount) error

func init() {
	ioc.Registry(
		NewAssociateTenantAccountWorkflow,
		workflows.NewAssociateTenantAccountWorkflow,
		tidbrepository.NewSaveTenantAccount,
		observability.NewObservability,
	)
}

func NewAssociateTenantAccountWorkflow(
	associateTenantAccountWorkflow workflows.AssociateTenantAccountWorkflow,
	saveTenantAccount tidbrepository.SaveTenantAccount,
	obs observability.Observability) AssociateTenantAccountWorkflow {
	return func(ctx context.Context, input domain.TenantAccount) error {
		// Obtener el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		workflow, err := associateTenantAccountWorkflow.Restore(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to restore workflow: %w", err)
		}
		if err := workflow.SetAssociationCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"email", input.Account.Email)
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = saveTenantAccount(ctx, input, fsmState)
		if err != nil {
			return fmt.Errorf("failed to save tenant account: %w", err)
		}
		return nil
	}
}
