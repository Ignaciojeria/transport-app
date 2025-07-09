package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/encryption"

	"github.com/biter777/countries"
	"github.com/google/uuid"
)

func MapClientCredentialsTable(e domain.ClientCredentials, encryptionService *encryption.ClientCredentialsEncryption) (table.ClientCredential, error) {
	// Encriptar el ClientSecret antes de guardarlo
	encryptedSecret, err := encryptionService.Encrypt(e.ClientSecret)
	if err != nil {
		return table.ClientCredential{}, err
	}

	return table.ClientCredential{
		ID:            e.ID,
		TenantID:      e.TenantID,
		ClientID:      e.ClientID,
		ClientSecret:  encryptedSecret,
		AllowedScopes: e.AllowedScopes,
		Status:        e.Status,
		CreatedAt:     e.CreatedAt,
		ExpiresAt:     e.ExpiresAt,
	}, nil
}

func MapClientCredentialsDomain(t table.ClientCredential, encryptionService *encryption.ClientCredentialsEncryption) (domain.ClientCredentials, error) {
	// Desencriptar el ClientSecret al leerlo
	decryptedSecret, err := encryptionService.Decrypt(t.ClientSecret)
	if err != nil {
		return domain.ClientCredentials{}, err
	}

	// Obtener el país del tenant desde el preload
	tenantCountry := countries.Unknown
	if t.Tenant.ID != uuid.Nil {
		tenantCountry = t.Tenant.Map().Country
	}

	return domain.ClientCredentials{
		ID:            t.ID,
		TenantID:      t.TenantID,
		TenantCountry: tenantCountry,
		ClientID:      t.ClientID,
		ClientSecret:  decryptedSecret,
		AllowedScopes: t.AllowedScopes,
		Status:        t.Status,
		CreatedAt:     t.CreatedAt,
		ExpiresAt:     t.ExpiresAt,
	}, nil
}
