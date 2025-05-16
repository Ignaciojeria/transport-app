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
		NewFindTenantByEmail,
		database.NewConnectionFactory,
	)
}

type FindTenantByEmail func(
	ctx context.Context,
	email string,
) (domain.Tenant, error)

func NewFindTenantByEmail(conn database.ConnectionFactory) FindTenantByEmail {
	return func(ctx context.Context, email string) (domain.Tenant, error) {
		// Crear una variable para almacenar la organización desde la tabla
		var tableOrg table.Tenant

		// Intentar buscar la organización por email
		err := conn.DB.WithContext(ctx).
			Where("email = ?", email).
			First(&tableOrg).Error

		// Manejar casos de error
		if err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.Tenant{}, errors.Wrap(ErrTenantNotFound, "tenant with email not found")
			}

			// Retornar cualquier otro error
			return domain.Tenant{}, err
		}

		return tableOrg.Map(), nil
	}
}
