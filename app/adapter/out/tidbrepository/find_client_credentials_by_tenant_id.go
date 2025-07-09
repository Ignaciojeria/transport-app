package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/infrastructure/encryption"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func init() {
	ioc.Registry(
		NewFindClientCredentialsByTenantID,
		database.NewConnectionFactory,
		encryption.NewClientCredentialsEncryptionFromConfig,
	)
}

type FindClientCredentialsByTenantID func(ctx context.Context, tenantID uuid.UUID) ([]domain.ClientCredentials, error)

func NewFindClientCredentialsByTenantID(conn database.ConnectionFactory, encryptionService *encryption.ClientCredentialsEncryption) FindClientCredentialsByTenantID {
	return func(ctx context.Context, tenantID uuid.UUID) ([]domain.ClientCredentials, error) {
		var credentials []table.ClientCredential
		err := conn.DB.WithContext(ctx).
			Table("client_credentials").
			Where("tenant_id = ?", tenantID).
			Find(&credentials).Error

		if err != nil {
			return nil, errors.Wrap(ErrClientCredentialsDatabase, "failed to find client credentials by tenant")
		}

		result := make([]domain.ClientCredentials, len(credentials))
		for i, cred := range credentials {
			domainCred, err := mapper.MapClientCredentialsDomain(cred, encryptionService)
			if err != nil {
				return nil, errors.Wrap(ErrClientCredentialsDatabase, "failed to decrypt client secret")
			}
			result[i] = domainCred
		}

		return result, nil
	}
}
