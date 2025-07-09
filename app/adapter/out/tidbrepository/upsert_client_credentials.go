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
		NewUpsertClientCredentials,
		database.NewConnectionFactory,
		encryption.NewClientCredentialsEncryptionFromConfig,
	)
}

type UpsertClientCredentials func(ctx context.Context, credentials domain.ClientCredentials) (domain.ClientCredentials, error)

func NewUpsertClientCredentials(conn database.ConnectionFactory, encryptionService *encryption.ClientCredentialsEncryption) UpsertClientCredentials {
	return func(ctx context.Context, credentials domain.ClientCredentials) (domain.ClientCredentials, error) {
		var existing table.ClientCredential
		err := conn.DB.WithContext(ctx).
			Table("client_credentials").
			Where("id = ?", credentials.ID).
			First(&existing).Error

		if err == nil {
			// Actualizar registro existente
			tableCredentials, err := mapper.MapClientCredentialsTable(credentials, encryptionService)
			if err != nil {
				return domain.ClientCredentials{}, errors.Wrap(ErrClientCredentialsDatabase, "failed to encrypt client secret")
			}
			if err := conn.DB.WithContext(ctx).Save(&tableCredentials).Error; err != nil {
				return domain.ClientCredentials{}, errors.Wrap(ErrClientCredentialsDatabase, "failed to update client credentials")
			}
			return mapper.MapClientCredentialsDomain(tableCredentials, encryptionService)
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ClientCredentials{}, err
		}

		// Crear nuevo registro
		tableCredentials, err := mapper.MapClientCredentialsTable(credentials, encryptionService)
		if err != nil {
			return domain.ClientCredentials{}, errors.Wrap(ErrClientCredentialsDatabase, "failed to encrypt client secret")
		}
		if err := conn.DB.WithContext(ctx).Create(&tableCredentials).Error; err != nil {
			return domain.ClientCredentials{}, errors.Wrap(ErrClientCredentialsDatabase, "failed to create client credentials")
		}

		return mapper.MapClientCredentialsDomain(tableCredentials, encryptionService)
	}
}
