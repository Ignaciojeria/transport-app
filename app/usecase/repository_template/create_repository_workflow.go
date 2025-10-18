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
	"github.com/go-git/go-git/v6/plumbing/object"
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
		// Usar /tmp para todos los casos (desarrollo y producción)
		repoPath := filepath.Join("/tmp", "tenants", repoName)

		// Intentar crear repositorio remoto si hay token configurado
		var repo *gitv6.Repository

		obs.Logger.InfoContext(ctx, "Checking Git configuration",
			"has_token", conf.GIT_TOKEN != "",
			"token_length", len(conf.GIT_TOKEN),
			"repository_path", repoPath)

		if conf.GIT_TOKEN != "" {
			// Crear repositorio en GitHub en la organización transport-app-agents
			obs.Logger.InfoContext(ctx, "Attempting to create public GitHub repository in organization",
				"repository_name", repoName,
				"organization", "transport-app-agents")
			githubRepo, err := gitrepository.CreateGitHubOrganizationRepository(ctx, conf.GIT_TOKEN, "transport-app-agents", repoName,
				fmt.Sprintf("Tenant repository for %s", repoName), false)
			if err != nil {
				obs.Logger.WarnContext(ctx, "Failed to create GitHub repository, falling back to local",
					"error", err, "repository_name", repoName)
				// Fallback a repositorio local
				repo, err = gitrepository.CreateRepositoryIfNotExists(repoPath)
			} else {
				obs.Logger.InfoContext(ctx, "Public GitHub repository created successfully in organization",
					"repository_name", repoName,
					"organization", "transport-app-agents",
					"url", githubRepo.HTMLURL)

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
		obs.Logger.InfoContext(ctx, "Checking gitRepoAdapter availability",
			"gitRepoAdapter_nil", gitRepoAdapter == nil,
			"repository_name", repoName)

		if gitRepoAdapter != nil {
			obs.Logger.InfoContext(ctx, "Starting template file creation process",
				"repository_name", repoName,
				"repository_path", repoPath)

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
				Author: &object.Signature{
					Name:  "Transport App",
					Email: "noreply@transport-app.com",
					When:  time.Now(),
				},
				Committer: &object.Signature{
					Name:  "Transport App",
					Email: "noreply@transport-app.com",
					When:  time.Now(),
				},
			}
			if err := gitRepoAdapter.CommitRepository(repoPath, "feat: initial repository setup for tenant", commitOptions); err != nil {
				obs.Logger.WarnContext(ctx, "Failed to create initial commit",
					"error", err,
					"repository_name", repoName,
					"repository_path", repoPath)
			} else {
				obs.Logger.InfoContext(ctx, "Initial commit created successfully",
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
				Force:      true,           // Force push para resolver conflictos
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
		} else {
			obs.Logger.WarnContext(ctx, "gitRepoAdapter is nil, skipping template file creation and push",
				"repository_name", repoName,
				"repository_path", repoPath)
		}

		// Log del repositorio creado (para debugging)
		_ = repo // Evitar warning de variable no usada

		return nil
	}
}

// createInitialTenantFiles crea archivos iniciales para el tenant usando templates embedded
func createInitialTenantFiles(repoPath, repoName string) error {
	fmt.Printf("DEBUG: Starting createInitialTenantFiles for %s at %s\n", repoName, repoPath)

	// Crear directorios necesarios
	dirs := []string{
		filepath.Join(repoPath, "data"),
		filepath.Join(repoPath, "config"),
		filepath.Join(repoPath, "assets"),
		filepath.Join(repoPath, ".github", "workflows"),
	}

	fmt.Printf("DEBUG: Creating directories: %v\n", dirs)
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("DEBUG: Failed to create directory %s: %v\n", dir, err)
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		fmt.Printf("DEBUG: Successfully created directory: %s\n", dir)
	}

	// Extraer tenant ID del nombre del repositorio
	tenantID := strings.TrimPrefix(repoName, "tenant-")

	// Crear README.md usando template
	fmt.Printf("DEBUG: Creating README.md\n")
	readmeContent := processTemplate(DeployReadmeTemplate, map[string]string{
		"TenantName": repoName,
		"TenantID":   tenantID,
		"CreatedAt":  time.Now().Format(time.RFC3339),
	})

	readmePath := filepath.Join(repoPath, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		fmt.Printf("DEBUG: Failed to create README.md: %v\n", err)
		return fmt.Errorf("failed to create README.md: %w", err)
	}
	fmt.Printf("DEBUG: Successfully created README.md at %s\n", readmePath)

	// Crear config/tenant.json usando template
	fmt.Printf("DEBUG: Creating tenant.json\n")
	configContent := processTemplate(TenantConfigTemplate, map[string]string{
		"TenantName": repoName,
		"TenantID":   tenantID,
		"CreatedAt":  time.Now().Format(time.RFC3339),
	})

	configPath := filepath.Join(repoPath, "config", "tenant.json")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		fmt.Printf("DEBUG: Failed to create tenant.json: %v\n", err)
		return fmt.Errorf("failed to create tenant.json: %w", err)
	}
	fmt.Printf("DEBUG: Successfully created tenant.json at %s\n", configPath)

	// Crear .gitignore usando template
	fmt.Printf("DEBUG: Creating .gitignore\n")
	gitignorePath := filepath.Join(repoPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(GitignoreTemplate), 0644); err != nil {
		fmt.Printf("DEBUG: Failed to create .gitignore: %v\n", err)
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}
	fmt.Printf("DEBUG: Successfully created .gitignore at %s\n", gitignorePath)

	// Crear index.html usando template
	fmt.Printf("DEBUG: Creating index.html\n")
	indexContent := processTemplate(IndexTemplate, map[string]string{
		"TenantName": repoName,
		"TenantID":   tenantID,
		"CreatedAt":  time.Now().Format(time.RFC3339),
	})

	indexPath := filepath.Join(repoPath, "index.html")
	if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
		fmt.Printf("DEBUG: Failed to create index.html: %v\n", err)
		return fmt.Errorf("failed to create index.html: %w", err)
	}
	fmt.Printf("DEBUG: Successfully created index.html at %s\n", indexPath)

	// Crear workflow de GitHub Actions usando template
	fmt.Printf("DEBUG: Creating GitHub workflow\n")
	workflowContent := processTemplate(GitHubWorkflowTemplate, map[string]string{
		"TenantName": repoName,
		"TenantID":   tenantID,
	})

	workflowPath := filepath.Join(repoPath, ".github", "workflows", "deploy.yml")
	if err := os.WriteFile(workflowPath, []byte(workflowContent), 0644); err != nil {
		fmt.Printf("DEBUG: Failed to create GitHub workflow: %v\n", err)
		return fmt.Errorf("failed to create GitHub workflow: %w", err)
	}
	fmt.Printf("DEBUG: Successfully created GitHub workflow at %s\n", workflowPath)

	// Crear firebase.json usando template
	fmt.Printf("DEBUG: Creating firebase.json\n")
	firebasePath := filepath.Join(repoPath, "firebase.json")
	if err := os.WriteFile(firebasePath, []byte(FirebaseTemplate), 0644); err != nil {
		fmt.Printf("DEBUG: Failed to create firebase.json: %v\n", err)
		return fmt.Errorf("failed to create firebase.json: %w", err)
	}
	fmt.Printf("DEBUG: Successfully created firebase.json at %s\n", firebasePath)

	fmt.Printf("DEBUG: All template files created successfully\n")
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
