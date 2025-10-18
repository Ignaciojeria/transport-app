package gitrepository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GitHubRepositoryRequest representa la estructura para crear un repositorio en GitHub
type GitHubRepositoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
}

// GitHubRepositoryResponse representa la respuesta de GitHub al crear un repositorio
type GitHubRepositoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	HTMLURL     string `json:"html_url"`
	CloneURL    string `json:"clone_url"`
	SSHURL      string `json:"ssh_url"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateGitHubRepository crea un repositorio en GitHub usando la API
func CreateGitHubRepository(ctx context.Context, token, repoName, description string, isPrivate bool) (*GitHubRepositoryResponse, error) {
	// Construir la URL de la API de GitHub
	url := "https://api.github.com/user/repos"

	// Crear la estructura de datos
	reqData := GitHubRepositoryRequest{
		Name:        repoName,
		Description: description,
		Private:     isPrivate,
		AutoInit:    false, // No crear README inicial, usaremos nuestro template
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Crear la petición HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Configurar headers
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de respuesta
	if resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("failed to create repository: HTTP %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("failed to create repository: %v", errorResp)
	}

	// Decodificar la respuesta
	var repoResp GitHubRepositoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&repoResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &repoResp, nil
}

// CreateGitLabRepository crea un repositorio en GitLab usando la API
func CreateGitLabRepository(ctx context.Context, token, repoName, description string, isPrivate bool) (*GitHubRepositoryResponse, error) {
	// Construir la URL de la API de GitLab
	url := "https://gitlab.com/api/v4/projects"

	// Crear la estructura de datos para GitLab
	reqData := map[string]interface{}{
		"name":                   repoName,
		"description":            description,
		"visibility":             map[string]bool{"private": isPrivate},
		"initialize_with_readme": false, // No crear README inicial, usaremos nuestro template
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	// Crear la petición HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Configurar headers
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Verificar el código de respuesta
	if resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("failed to create repository: HTTP %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("failed to create repository: %v", errorResp)
	}

	// Decodificar la respuesta de GitLab
	var gitlabResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&gitlabResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convertir a la estructura estándar
	repoResp := &GitHubRepositoryResponse{
		Name:        gitlabResp["name"].(string),
		Description: gitlabResp["description"].(string),
		Private:     gitlabResp["visibility"].(string) == "private",
		HTMLURL:     gitlabResp["web_url"].(string),
		CloneURL:    gitlabResp["http_url_to_repo"].(string),
		SSHURL:      gitlabResp["ssh_url_to_repo"].(string),
		CreatedAt:   gitlabResp["created_at"].(string),
		UpdatedAt:   gitlabResp["updated_at"].(string),
	}

	return repoResp, nil
}
