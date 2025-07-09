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
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewFindClientCredentialsByID,
		database.NewConnectionFactory,
		encryption.NewClientCredentialsEncryptionFromConfig,
	)
}

type FindClientCredentialsByID func(ctx context.Context, id uuid.UUID) (domain.ClientCredentials, error)

func NewFindClientCredentialsByID(conn database.ConnectionFactory, encryptionService *encryption.ClientCredentialsEncryption) FindClientCredentialsByID {
	return func(ctx context.Context, id uuid.UUID) (domain.ClientCredentials, error) {
		var credentials table.ClientCredential
		err := conn.DB.WithContext(ctx).
			Table("client_credentials").
			Where("id = ?", id).
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
