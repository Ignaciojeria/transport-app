package domain

import (
	"time"

	"github.com/google/uuid"
)

type ClientCredentials struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	ClientID      string
	ClientSecret  string
	AllowedScopes []string
	Status        string
	CreatedAt     time.Time
	ExpiresAt     *time.Time
}

// NewClientCredentials crea una nueva instancia de ClientCredentials
func NewClientCredentials(
	tenantID uuid.UUID,
	clientID string,
	clientSecret string,
	allowedScopes []string,
) *ClientCredentials {
	return &ClientCredentials{
		ID:            uuid.New(),
		TenantID:      tenantID,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		AllowedScopes: allowedScopes,
		Status:        "active",
		CreatedAt:     time.Now(),
	}
}

// IsActive verifica si las credenciales están activas
func (c *ClientCredentials) IsActive() bool {
	return c.Status == "active"
}

// IsExpired verifica si las credenciales han expirado
func (c *ClientCredentials) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// IsValid verifica si las credenciales son válidas (activas y no expiradas)
func (c *ClientCredentials) IsValid() bool {
	return c.IsActive() && !c.IsExpired()
}

// HasScope verifica si las credenciales tienen un scope específico
func (c *ClientCredentials) HasScope(scope string) bool {
	for _, s := range c.AllowedScopes {
		if s == scope {
			return true
		}
	}
	return false
}
