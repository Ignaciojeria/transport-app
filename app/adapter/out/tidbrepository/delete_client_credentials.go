package tidbrepository

import (
	"context"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

func init() {
	ioc.Registry(
		NewDeleteClientCredentials,
		database.NewConnectionFactory,
	)
}

type DeleteClientCredentials func(ctx context.Context, id uuid.UUID) error

func NewDeleteClientCredentials(conn database.ConnectionFactory) DeleteClientCredentials {
	return func(ctx context.Context, id uuid.UUID) error {
		result := conn.DB.WithContext(ctx).
			Table("client_credentials").
			Where("id = ?", id).
			Delete(&struct{}{})

		if result.Error != nil {
			return errors.Wrap(ErrClientCredentialsDatabase, "failed to delete client credentials")
		}

		if result.RowsAffected == 0 {
			return errors.Wrap(ErrClientCredentialsNotFound, "client credentials not found for deletion")
		}

		return nil
	}
}
