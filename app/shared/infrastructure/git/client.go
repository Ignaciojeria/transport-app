package git

import (
	"fmt"
	"os"
	"transport-app/app/shared/configuration"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/go-git/go-git/v6"
)

func init() {
	ioc.Registry(NewClient, configuration.NewConf)
}

// GitClient define la interfaz para operaciones de Git
type GitClient interface {
	GetRepository() *git.Repository
	GetRepositoryPath() string
}

// gitClient implementa la interfaz GitClient
type gitClient struct {
	repo *git.Repository
	path string
}

// NewClient crea una nueva instancia del cliente Git
func NewClient(conf configuration.Conf) (GitClient, error) {
	// Si no se especifica un repositorio, retornamos nil
	if conf.GIT_REPOSITORY_PATH == "" {
		fmt.Println("Git client will not be initialized because GIT_REPOSITORY_PATH is not set")
		return nil, nil
	}

	// Crear el directorio si no existe
	if _, err := os.Stat(conf.GIT_REPOSITORY_PATH); os.IsNotExist(err) {
		fmt.Printf("Creating agent repository directory: %s\n", conf.GIT_REPOSITORY_PATH)
		if err := os.MkdirAll(conf.GIT_REPOSITORY_PATH, 0755); err != nil {
			return nil, fmt.Errorf("failed to create agent repository directory: %w", err)
		}
	}

	// Verificar si ya es un repositorio Git
	repo, err := git.PlainOpen(conf.GIT_REPOSITORY_PATH)
	if err != nil {
		// Si no es un repositorio Git, inicializarlo
		fmt.Printf("Initializing Git repository in: %s\n", conf.GIT_REPOSITORY_PATH)
		repo, err = git.PlainInit(conf.GIT_REPOSITORY_PATH, false)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize git repository: %w", err)
		}
	}

	return &gitClient{
		repo: repo,
		path: conf.GIT_REPOSITORY_PATH,
	}, nil
}

// GetRepository retorna la instancia del repositorio Git
func (gc *gitClient) GetRepository() *git.Repository {
	return gc.repo
}

// GetRepositoryPath retorna la ruta del repositorio
func (gc *gitClient) GetRepositoryPath() string {
	return gc.path
}

// CreateRepositoryIfNotExists crea un repositorio Git si no existe
func CreateRepositoryIfNotExists(path string) (*git.Repository, error) {
	// Verificar si la ruta existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Crear el directorio si no existe
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Verificar si ya es un repositorio Git
	if _, err := git.PlainOpen(path); err == nil {
		// Ya es un repositorio, abrirlo
		return git.PlainOpen(path)
	}

	// Crear un nuevo repositorio
	repo, err := git.PlainInit(path, false)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize git repository: %w", err)
	}

	fmt.Printf("Git repository created at: %s\n", path)
	return repo, nil
}

// IsGitRepository verifica si la ruta es un repositorio Git válido
func IsGitRepository(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}
