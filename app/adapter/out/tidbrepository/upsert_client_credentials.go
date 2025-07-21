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
		NewSaveFSMTransition,
	)
}

type UpsertClientCredentials func(ctx context.Context, credentials domain.ClientCredentials, fsmState ...domain.FSMState) (domain.ClientCredentials, error)

func NewUpsertClientCredentials(conn database.ConnectionFactory, encryptionService *encryption.ClientCredentialsEncryption, saveFSMTransition SaveFSMTransition) UpsertClientCredentials {
	return func(ctx context.Context, credentials domain.ClientCredentials, fsmState ...domain.FSMState) (domain.ClientCredentials, error) {
		var result domain.ClientCredentials

		err := conn.Transaction(func(tx *gorm.DB) error {
			var existing table.ClientCredential
			err := tx.WithContext(ctx).
				Table("client_credentials").
				Where("id = ?", credentials.ID).
				First(&existing).Error

			var tableCredentials table.ClientCredential
			if err == nil {
				// Actualizar registro existente
				tableCredentials, err = mapper.MapClientCredentialsTable(credentials, encryptionService)
				if err != nil {
					return errors.Wrap(ErrClientCredentialsDatabase, "failed to encrypt client secret")
				}
				if err := tx.WithContext(ctx).Save(&tableCredentials).Error; err != nil {
					return errors.Wrap(ErrClientCredentialsDatabase, "failed to update client credentials")
				}
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				// Crear nuevo registro
				tableCredentials, err = mapper.MapClientCredentialsTable(credentials, encryptionService)
				if err != nil {
					return errors.Wrap(ErrClientCredentialsDatabase, "failed to encrypt client secret")
				}
				if err := tx.WithContext(ctx).Create(&tableCredentials).Error; err != nil {
					return errors.Wrap(ErrClientCredentialsDatabase, "failed to create client credentials")
				}
			} else {
				return err
			}

			// Mapear resultado
			var mapErr error
			result, mapErr = mapper.MapClientCredentialsDomain(tableCredentials, encryptionService)
			if mapErr != nil {
				return errors.Wrap(ErrClientCredentialsDatabase, "failed to map client credentials")
			}

			// Persistir FSMState si está presente
			if len(fsmState) > 0 {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})

		if err != nil {
			return domain.ClientCredentials{}, err
		}

		return result, nil
	}
}
