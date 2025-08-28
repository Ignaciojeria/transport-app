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

type CreateDefaultClientCredentialsWorkflow func(ctx context.Context, input domain.ClientCredentials) error

func init() {
	ioc.Registry(
		NewCreateDefaultClientCredentialsWorkflow,
		workflows.NewGenericWorkflow,
		tidbrepository.NewUpsertClientCredentials,
		observability.NewObservability,
	)
}

func NewCreateDefaultClientCredentialsWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	upsertClientCredentials tidbrepository.UpsertClientCredentials,
	obs observability.Observability,
) CreateDefaultClientCredentialsWorkflow {
	return func(ctx context.Context, input domain.ClientCredentials) error {
		// Obtener el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}
		
		// Configurar workflow genérico para client credentials creation
		config := workflows.CreateWorkflow("client_credentials", "create")
		
		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}

		// Intentar transición a credenciales creadas
		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"client_id", input.ClientID)
			return nil
		}

		// Mapear el estado FSM
		fsmState := workflow.Map(ctx)

		// Guardar las credenciales con el estado FSM
		_, err = upsertClientCredentials(ctx, input, fsmState)
		if err != nil {
			return fmt.Errorf("failed to upsert client credentials: %w", err)
		}

		return nil
	}
}
