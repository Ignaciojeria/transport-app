package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewFindOrganizationByEmail,
		tidb.NewTIDBConnection,
	)
}

type FindOrganizationByEmail func(
	ctx context.Context,
	email string,
) (domain.Organization, error)

func NewFindOrganizationByEmail(conn tidb.TIDBConnection) FindOrganizationByEmail {
	return func(ctx context.Context, email string) (domain.Organization, error) {
		// Crear una variable para almacenar la organización desde la tabla
		var tableOrg table.Organization

		// Intentar buscar la organización por email
		err := conn.DB.WithContext(ctx).
			Where("email = ?", email).
			First(&tableOrg).Error

		// Manejar casos de error
		if err != nil {
			// Si no se encuentra, retornar un error específico
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.Organization{}, ErrOrganizationNotFound.
					New("organization with email not found: %v", err)
			}

			// Retornar cualquier otro error
			return domain.Organization{}, err
		}

		return tableOrg.Map(), nil
	}
}
