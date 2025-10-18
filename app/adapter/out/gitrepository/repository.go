package gitrepository

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"transport-app/app/shared/infrastructure/git"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	gitv6 "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
)

func init() {
	ioc.Registry(NewGitRepositoryAdapter, git.NewClient)
}

// GitRepositoryAdapter proporciona implementaciones específicas para operaciones de Git
type GitRepositoryAdapter struct {
	client git.GitClient
}

// NewGitRepositoryAdapter crea una nueva instancia del adaptador Git
func NewGitRepositoryAdapter(client git.GitClient) *GitRepositoryAdapter {
	return &GitRepositoryAdapter{
		client: client,
	}
}

// CreateRepositoryIfNotExists crea un repositorio Git en la ruta especificada si no existe
func CreateRepositoryIfNotExists(path string) (*gitv6.Repository, error) {
	// Verificar si la ruta existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Crear el directorio si no existe
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Verificar si ya es un repositorio Git
	if _, err := gitv6.PlainOpen(path); err == nil {
		// Ya es un repositorio, abrirlo
		return gitv6.PlainOpen(path)
	}

	// Crear un nuevo repositorio
	repo, err := gitv6.PlainInit(path, false)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize git repository: %w", err)
	}

	fmt.Printf("Git repository created at: %s\n", path)
	return repo, nil
}

// PushOptions define las opciones para hacer push
type PushOptions struct {
	Token      string
	SSHKeyPath string
	RemoteName string
	Progress   io.Writer
	Prune      bool
	Force      bool
	FollowTags bool
	Atomic     bool
}

// PullOptions define las opciones para hacer pull
type PullOptions struct {
	Token         string
	SSHKeyPath    string
	RemoteName    string
	ReferenceName plumbing.ReferenceName
	SingleBranch  bool
	Depth         int
	Progress      io.Writer
}

// CommitOptions define las opciones para crear un commit
type CommitOptions struct {
	All       bool
	Author    *object.Signature
	Committer *object.Signature
}

// Push envía los cambios al repositorio remoto
func (gra *GitRepositoryAdapter) Push(ctx context.Context, options *PushOptions) error {
	if gra.client == nil {
		return fmt.Errorf("git client is not initialized")
	}

	repo := gra.client.GetRepository()
	if repo == nil {
		return fmt.Errorf("repository is not available")
	}

	if options == nil {
		options = &PushOptions{}
	}

	// Configurar autenticación (priorizar SSH en Gitpod)
	var auth transport.AuthMethod
	if options.SSHKeyPath != "" {
		sshKey, err := ssh.NewPublicKeysFromFile("git", options.SSHKeyPath, "")
		if err != nil {
			return fmt.Errorf("failed to load SSH key: %w", err)
		}
		auth = sshKey
	} else if options.Token != "" {
		auth = &http.BasicAuth{
			Username: "token", // Para GitHub, GitLab, etc.
			Password: options.Token,
		}
	}

	return repo.PushContext(ctx, &gitv6.PushOptions{
		RemoteName: options.RemoteName,
		Auth:       auth,
		Progress:   options.Progress,
		Prune:      options.Prune,
		Force:      options.Force,
		FollowTags: options.FollowTags,
		Atomic:     options.Atomic,
	})
}

// Pull actualiza el repositorio con los cambios remotos
func (gra *GitRepositoryAdapter) Pull(ctx context.Context, options *PullOptions) error {
	if gra.client == nil {
		return fmt.Errorf("git client is not initialized")
	}

	repo := gra.client.GetRepository()
	if repo == nil {
		return fmt.Errorf("repository is not available")
	}

	if options == nil {
		options = &PullOptions{}
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Configurar autenticación (priorizar SSH en Gitpod)
	var auth transport.AuthMethod
	if options.SSHKeyPath != "" {
		sshKey, err := ssh.NewPublicKeysFromFile("git", options.SSHKeyPath, "")
		if err != nil {
			return fmt.Errorf("failed to load SSH key: %w", err)
		}
		auth = sshKey
	} else if options.Token != "" {
		auth = &http.BasicAuth{
			Username: "token", // Para GitHub, GitLab, etc.
			Password: options.Token,
		}
	}

	return workTree.PullContext(ctx, &gitv6.PullOptions{
		RemoteName:    options.RemoteName,
		ReferenceName: options.ReferenceName,
		SingleBranch:  options.SingleBranch,
		Depth:         options.Depth,
		Auth:          auth,
		Progress:      options.Progress,
	})
}

// AddStage agrega archivos al staging area
func (gra *GitRepositoryAdapter) AddStage(path string) error {
	if gra.client == nil {
		return fmt.Errorf("git client is not initialized")
	}

	repo := gra.client.GetRepository()
	if repo == nil {
		return fmt.Errorf("repository is not available")
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	_, err = workTree.Add(path)
	return err
}

// Commit crea un nuevo commit
func (gra *GitRepositoryAdapter) Commit(message string, options *CommitOptions) error {
	if gra.client == nil {
		return fmt.Errorf("git client is not initialized")
	}

	repo := gra.client.GetRepository()
	if repo == nil {
		return fmt.Errorf("repository is not available")
	}

	if options == nil {
		options = &CommitOptions{}
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	_, err = workTree.Commit(message, &gitv6.CommitOptions{
		All:       options.All,
		Author:    options.Author,
		Committer: options.Committer,
	})

	return err
}

// GetStatus obtiene el estado del repositorio
func (gra *GitRepositoryAdapter) GetStatus() (gitv6.Status, error) {
	if gra.client == nil {
		return nil, fmt.Errorf("git client is not initialized")
	}

	repo := gra.client.GetRepository()
	if repo == nil {
		return nil, fmt.Errorf("repository is not available")
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	return workTree.Status()
}

// GetCurrentBranch obtiene la rama actual
func (gra *GitRepositoryAdapter) GetCurrentBranch() (string, error) {
	if gra.client == nil {
		return "", fmt.Errorf("git client is not initialized")
	}

	repo := gra.client.GetRepository()
	if repo == nil {
		return "", fmt.Errorf("repository is not available")
	}

	ref, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return ref.Name().Short(), nil
}

// GetCurrentCommit obtiene el hash del commit actual
func (gra *GitRepositoryAdapter) GetCurrentCommit() (string, error) {
	if gra.client == nil {
		return "", fmt.Errorf("git client is not initialized")
	}

	repo := gra.client.GetRepository()
	if repo == nil {
		return "", fmt.Errorf("repository is not available")
	}

	ref, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get current commit: %w", err)
	}
	return ref.Hash().String(), nil
}

// IsDirty verifica si el repositorio tiene cambios sin commitear
func (gra *GitRepositoryAdapter) IsDirty() (bool, error) {
	if gra.client == nil {
		return false, fmt.Errorf("git client is not initialized")
	}

	repo := gra.client.GetRepository()
	if repo == nil {
		return false, fmt.Errorf("repository is not available")
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := workTree.Status()
	if err != nil {
		return false, fmt.Errorf("failed to get status: %w", err)
	}

	return !status.IsClean(), nil
}

// GetRelativePath obtiene la ruta relativa de un archivo respecto al repositorio
func (gra *GitRepositoryAdapter) GetRelativePath(filePath string) (string, error) {
	if gra.client == nil {
		return "", fmt.Errorf("git client is not initialized")
	}

	absRepoPath, err := filepath.Abs(gra.client.GetRepositoryPath())
	if err != nil {
		return "", fmt.Errorf("failed to get absolute repository path: %w", err)
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute file path: %w", err)
	}

	relPath, err := filepath.Rel(absRepoPath, absFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}

	return relPath, nil
}
