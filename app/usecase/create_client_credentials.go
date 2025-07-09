package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
)

type CreateClientCredentials func(ctx context.Context, tenantID uuid.UUID, scopes []string) (domain.ClientCredentials, error)

func init() {
	ioc.Registry(
		NewCreateClientCredentials,
		tidbrepository.NewUpsertClientCredentials,
	)
}

func NewCreateClientCredentials(upsertClientCredentials tidbrepository.UpsertClientCredentials) CreateClientCredentials {
	return func(ctx context.Context, tenantID uuid.UUID, scopes []string) (domain.ClientCredentials, error) {
		// Generar ClientID único
		clientID, err := generateClientID()
		if err != nil {
			return domain.ClientCredentials{}, fmt.Errorf("error generando client ID: %w", err)
		}

		// Generar ClientSecret seguro
		clientSecret, err := generateClientSecret()
		if err != nil {
			return domain.ClientCredentials{}, fmt.Errorf("error generando client secret: %w", err)
		}

		// Crear las credenciales
		credentials := domain.NewClientCredentials(
			tenantID,
			clientID,
			clientSecret,
			scopes,
		)

		// Guardar en la base de datos (se encripta automáticamente)
		savedCredentials, err := upsertClientCredentials(ctx, *credentials)
		if err != nil {
			return domain.ClientCredentials{}, fmt.Errorf("error guardando client credentials: %w", err)
		}

		return savedCredentials, nil
	}
}

// generateClientID genera un ClientID único
func generateClientID() (string, error) {
	// Generar 16 bytes aleatorios
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Codificar en base64 y limpiar caracteres especiales
	clientID := base64.URLEncoding.EncodeToString(bytes)
	// Remover padding y caracteres especiales para hacerlo más legible
	clientID = clientID[:24] // Tomar solo 24 caracteres para hacerlo más corto

	return clientID, nil
}

// generateClientSecret genera un ClientSecret seguro
func generateClientSecret() (string, error) {
	// Generar 32 bytes aleatorios para un secret fuerte
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Codificar en base64
	clientSecret := base64.URLEncoding.EncodeToString(bytes)
	// Remover padding para hacerlo más legible
	clientSecret = clientSecret[:43] // Tomar solo 43 caracteres

	return clientSecret, nil
}
