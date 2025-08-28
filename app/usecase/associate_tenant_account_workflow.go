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
		workflows.NewGenericWorkflow,
		tidbrepository.NewSaveTenantAccount,
		observability.NewObservability,
	)
}

func NewAssociateTenantAccountWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	saveTenantAccount tidbrepository.SaveTenantAccount,
	obs observability.Observability) AssociateTenantAccountWorkflow {
	return func(ctx context.Context, input domain.TenantAccount) error {
		// Obtener el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		
		// Configurar workflow genérico para tenant-account association
		config := workflows.WorkflowConfig{
			Name:           "associate_tenant_account_workflow",
			StartedState:   "association_started",
			CompletedState: "association_completed",
			UseStorjBucket: false,
		}
		
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
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
