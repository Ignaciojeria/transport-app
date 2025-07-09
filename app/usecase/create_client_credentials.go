package usecase

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"github.com/google/uuid"
)

type CreateClientCredentials func(ctx context.Context, tenantID uuid.UUID, tenantCountry countries.CountryCode, scopes []string) (domain.ClientCredentials, error)

func init() {
	ioc.Registry(
		NewCreateClientCredentials,
		tidbrepository.NewUpsertClientCredentials,
	)
}

func NewCreateClientCredentials(upsertClientCredentials tidbrepository.UpsertClientCredentials) CreateClientCredentials {
	return func(ctx context.Context, tenantID uuid.UUID, tenantCountry countries.CountryCode, scopes []string) (domain.ClientCredentials, error) {
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
			tenantCountry,
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

// generateClientID genera un ClientID único usando UUID
func generateClientID() (string, error) {
	// Generar un UUID v4 para el client ID
	clientID := uuid.New()
	return clientID.String(), nil
}

// generateClientSecret genera un ClientSecret seguro usando UUID
func generateClientSecret() (string, error) {
	// Generar un UUID v4 para el client secret
	clientSecret := uuid.New()
	return clientSecret.String(), nil
}
