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
)

func init() {
	ioc.Registry(
		NewFindActiveClientCredentials,
		database.NewConnectionFactory,
		encryption.NewClientCredentialsEncryptionFromConfig,
	)
}

type FindActiveClientCredentials func(ctx context.Context) ([]domain.ClientCredentials, error)

func NewFindActiveClientCredentials(conn database.ConnectionFactory, encryptionService *encryption.ClientCredentialsEncryption) FindActiveClientCredentials {
	return func(ctx context.Context) ([]domain.ClientCredentials, error) {
		var credentials []table.ClientCredential
		err := conn.DB.WithContext(ctx).
			Table("client_credentials").
			Where("status = ?", "active").
			Find(&credentials).Error

		if err != nil {
			return nil, errors.Wrap(ErrClientCredentialsDatabase, "failed to find active client credentials")
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
