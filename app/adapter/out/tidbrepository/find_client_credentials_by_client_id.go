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
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewFindClientCredentialsByClientID,
		database.NewConnectionFactory,
		encryption.NewClientCredentialsEncryptionFromConfig,
	)
}

type FindClientCredentialsByClientID func(ctx context.Context, clientID string) (domain.ClientCredentials, error)

func NewFindClientCredentialsByClientID(conn database.ConnectionFactory, encryptionService *encryption.ClientCredentialsEncryption) FindClientCredentialsByClientID {
	return func(ctx context.Context, clientID string) (domain.ClientCredentials, error) {
		var credentials table.ClientCredential
		err := conn.DB.WithContext(ctx).
			Preload("Tenant").
			Table("client_credentials").
			Where("client_id = ?", clientID).
			First(&credentials).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ClientCredentials{}, errors.Wrap(ErrClientCredentialsNotFound, "client credentials not found")
			}
			return domain.ClientCredentials{}, errors.Wrap(ErrClientCredentialsDatabase, "failed to find client credentials")
		}

		return mapper.MapClientCredentialsDomain(credentials, encryptionService)
	}
}
