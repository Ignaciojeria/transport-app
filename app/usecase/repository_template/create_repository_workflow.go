package usecase

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"transport-app/app/adapter/out/gitrepository"
	"transport-app/app/domain/workflows"
	"transport-app/app/shared/configuration"
	"transport-app/app/shared/infrastructure/observability"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	gitv6 "github.com/go-git/go-git/v6"
)

// Los templates están definidos en templates.go

type CreateRepositoryWorkflow func(ctx context.Context, repoName string) error

func init() {
	ioc.Registry(
		NewCreateRepositoryWorkflow,
		workflows.NewGenericWorkflow,
		gitrepository.NewGitRepositoryAdapter,
		observability.NewObservability,
		configuration.NewConf,
	)
}

func NewCreateRepositoryWorkflow(
	genericWorkflow workflows.GenericWorkflow,
	gitRepoAdapter *gitrepository.GitRepositoryAdapter,
	obs observability.Observability,
	conf configuration.Conf) CreateRepositoryWorkflow {
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
		repoPath := filepath.Join("./agent-repo", "tenants", repoName)

		// Intentar crear repositorio remoto si hay token configurado
		var repo *gitv6.Repository

		obs.Logger.InfoContext(ctx, "Checking Git configuration",
			"has_token", conf.GIT_TOKEN != "",
			"token_length", len(conf.GIT_TOKEN))

		if conf.GIT_TOKEN != "" {
			// Crear repositorio en GitHub
			obs.Logger.InfoContext(ctx, "Attempting to create GitHub repository",
				"repository_name", repoName)
			githubRepo, err := gitrepository.CreateGitHubRepository(ctx, conf.GIT_TOKEN, repoName,
				fmt.Sprintf("Tenant repository for %s", repoName), false)
			if err != nil {
				obs.Logger.WarnContext(ctx, "Failed to create GitHub repository, falling back to local",
					"error", err, "repository_name", repoName)
				// Fallback a repositorio local
				repo, err = gitrepository.CreateRepositoryIfNotExists(repoPath)
			} else {
				obs.Logger.InfoContext(ctx, "GitHub repository created successfully",
					"repository_name", repoName, "url", githubRepo.HTMLURL)

				// Limpiar el directorio si existe para poder clonar
				if err := os.RemoveAll(repoPath); err != nil {
					obs.Logger.WarnContext(ctx, "Failed to remove existing directory, falling back to local",
						"error", err, "repository_name", repoName)
					// Fallback a repositorio local
					repo, err = gitrepository.CreateRepositoryIfNotExists(repoPath)
				} else {
					// Clonar el repositorio remoto
					pushOptions := &gitrepository.PushOptions{
						Token: conf.GIT_TOKEN,
					}
					repo, err = gitRepoAdapter.CreateRemoteRepository(ctx, repoPath, repoName, githubRepo.CloneURL, pushOptions)
					if err != nil {
						obs.Logger.WarnContext(ctx, "Failed to clone remote repository, falling back to local",
							"error", err, "repository_name", repoName)
						// Fallback a repositorio local
						repo, err = gitrepository.CreateRepositoryIfNotExists(repoPath)
					}
				}
			}
		} else {
			// Crear repositorio local
			repo, err = gitrepository.CreateRepositoryIfNotExists(repoPath)
		}

		if err != nil {
			return fmt.Errorf("failed to create repository: %w", err)
		}

		obs.Logger.InfoContext(ctx,
			"Repository created successfully",
			"repository_name", repoName,
			"repository_path", repoPath,
			"workflow_state", fsmState)

		// Hacer push del repositorio si el adaptador está disponible
		if gitRepoAdapter != nil {
			// Crear archivos iniciales para el tenant
			if err := createInitialTenantFiles(repoPath, repoName); err != nil {
				obs.Logger.WarnContext(ctx, "Failed to create initial tenant files",
					"error", err,
					"repository_name", repoName,
					"repository_path", repoPath)
			} else {
				obs.Logger.InfoContext(ctx, "Initial tenant files created successfully",
					"repository_name", repoName,
					"repository_path", repoPath)
			}

			// Verificar el estado del repositorio antes del commit
			status, err := gitRepoAdapter.GetRepositoryStatus(repoPath)
			if err != nil {
				obs.Logger.WarnContext(ctx, "Failed to get repository status",
					"error", err,
					"repository_name", repoName)
			} else {
				obs.Logger.InfoContext(ctx, "Repository status before commit",
					"repository_name", repoName,
					"status", status)
			}

			// Crear un commit inicial en el repositorio específico
			commitOptions := &gitrepository.CommitOptions{
				All: true,
			}
			if err := gitRepoAdapter.CommitRepository(repoPath, "feat: initial repository setup for tenant", commitOptions); err != nil {
				obs.Logger.WarnContext(ctx, "Failed to create initial commit",
					"error", err,
					"repository_name", repoName,
					"repository_path", repoPath)
			}

			// Hacer push del repositorio específico
			// Usar token para Cloud Run, SSH para Gitpod
			obs.Logger.InfoContext(ctx, "Configuring push options",
				"has_token", conf.GIT_TOKEN != "",
				"token_length", len(conf.GIT_TOKEN))

			pushOptions := &gitrepository.PushOptions{
				Token:      conf.GIT_TOKEN, // Token para Cloud Run
				RemoteName: "origin",       // Nombre del remoto por defecto
				Progress:   nil,            // Sin progress para evitar logs excesivos
			}

			obs.Logger.InfoContext(ctx, "Attempting to push repository",
				"repository_name", repoName,
				"repository_path", repoPath)

			if err := gitRepoAdapter.PushRepository(ctx, repoPath, pushOptions); err != nil {
				// Si es un repositorio local sin remoto, es normal
				if strings.Contains(err.Error(), "no remote repositories configured") {
					obs.Logger.InfoContext(ctx, "Repository created locally (no remote configured)",
						"repository_name", repoName,
						"repository_path", repoPath)
				} else {
					obs.Logger.WarnContext(ctx, "Failed to push repository",
						"error", err,
						"repository_name", repoName,
						"repository_path", repoPath)
				}
				// No fallamos el workflow si el push falla
			} else {
				obs.Logger.InfoContext(ctx, "Repository pushed successfully",
					"repository_name", repoName,
					"repository_path", repoPath)
			}
		}

		// Log del repositorio creado (para debugging)
		_ = repo // Evitar warning de variable no usada

		return nil
	}
}

// createInitialTenantFiles crea archivos iniciales para el tenant usando templates embedded
func createInitialTenantFiles(repoPath, repoName string) error {
	// Crear directorios necesarios
	dirs := []string{
		filepath.Join(repoPath, "data"),
		filepath.Join(repoPath, "config"),
		filepath.Join(repoPath, "assets"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Extraer tenant ID del nombre del repositorio
	tenantID := strings.TrimPrefix(repoName, "tenant-")

	// Crear README.md usando template
	readmeContent := processTemplate(TenantRepoReadmeTemplate, map[string]string{
		"TenantName": repoName,
		"CreatedAt":  time.Now().Format(time.RFC3339),
	})

	readmePath := filepath.Join(repoPath, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// Crear config/tenant.json usando template
	configContent := processTemplate(TenantConfigTemplate, map[string]string{
		"TenantName": repoName,
		"TenantID":   tenantID,
		"CreatedAt":  time.Now().Format(time.RFC3339),
	})

	configPath := filepath.Join(repoPath, "config", "tenant.json")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create tenant.json: %w", err)
	}

	// Crear .gitignore usando template
	gitignorePath := filepath.Join(repoPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(GitignoreTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// processTemplate procesa un template simple reemplazando placeholders
func processTemplate(template string, data map[string]string) string {
	result := template
	for key, value := range data {
		placeholder := "{{." + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
