package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	"github.com/cockroachdb/errors"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewFindOrganizationByEmail,
		database.NewConnectionFactory,
	)
}

type FindOrganizationByEmail func(
	ctx context.Context,
	email string,
) (domain.Organization, error)

func NewFindOrganizationByEmail(conn database.ConnectionFactory) FindOrganizationByEmail {
	return func(ctx context.Context, email string) (domain.Organization, error) {
		// Crear una variable para almacenar la organización desde la tabla
		var tableOrg table.Organization

		// Intentar buscar la organización por email
		err := conn.DB.WithContext(ctx).
			Where("email = ?", email).
			First(&tableOrg).Error

		// Manejar casos de error
		if err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.Organization{}, errors.Wrap(ErrOrganizationNotFound, "organization with email not found")
			}

			// Retornar cualquier otro error
			return domain.Organization{}, err
		}

		return tableOrg.Map(), nil
	}
}
