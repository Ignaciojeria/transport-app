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

type CreateAccountWorkflow func(ctx context.Context, input domain.Account) error

func init() {
	ioc.Registry(
		NewCreateAccountWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertAccount,
		observability.NewObservability,
	)
}

func NewCreateAccountWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertAccount tidbrepository.UpsertAccount,
	obs observability.Observability) CreateAccountWorkflow {
	return func(ctx context.Context, input domain.Account) error {
		// Obtener el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		
		// Configurar workflow genérico para account creation
		config := workflows.CreateWorkflow("account", "create")
		
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"email", input.Email)
			return nil
		}
		fsmState := workflow.Map(ctx)
		err = upsertAccount(ctx, input, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert account: %w", err)
		}
		return nil
	}
}
