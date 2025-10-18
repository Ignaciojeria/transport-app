package usecase

import (
	"context"
	"fmt"
	"path/filepath"
	"transport-app/app/adapter/out/gitrepository"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateRepositoryWorkflow func(ctx context.Context, repoName string) error

func init() {
	ioc.Registry(
		NewCreateRepositoryWorkflow,
		workflows.NewGenericWorkflow,
		observability.NewObservability,
	)
}

func NewCreateRepositoryWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	obs observability.Observability) CreateRepositoryWorkflow {
	return func(ctx context.Context, repoName string) error {

		// Obtener el idempotency key desde el contexto
		key, ok := sharedcontext.IdempotencyKeyFromContext(ctx)
		if !ok {
			return fmt.Errorf("idempotency key not found in context")
		}

		// Configurar workflow genérico para repository creation
		config := workflows.CreateWorkflow("repository", "create")

		workflow, err := genericWorkflow.Initialize(ctx, key, config)
		if err != nil {
			return fmt.Errorf("failed to initialize workflow: %w", err)
		}

		if err := workflow.SetCompletedTransition(ctx); err != nil {
			obs.Logger.WarnContext(ctx,
				err.Error(),
				"repository_name", repoName)
			return nil
		}

		fsmState := workflow.Map(ctx)

		// Crear el repositorio usando la función independiente
		repoPath := filepath.Join("/tmp/repositories", repoName)
		repo, err := gitrepository.CreateRepositoryIfNotExists(repoPath)
		if err != nil {
			return fmt.Errorf("failed to create repository: %w", err)
		}

		obs.Logger.InfoContext(ctx,
			"Repository created successfully",
			"repository_name", repoName,
			"repository_path", repoPath,
			"workflow_state", fsmState)

		// Log del repositorio creado (para debugging)
		_ = repo // Evitar warning de variable no usada

		return nil
	}
}
