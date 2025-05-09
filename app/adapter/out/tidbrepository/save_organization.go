package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewSaveOrganization,
		database.NewConnectionFactory,
	)
}

type SaveOrganization func(
	context.Context,
	domain.Organization,
) (domain.Organization, error)

func NewSaveOrganization(conn database.ConnectionFactory) SaveOrganization {
	return func(ctx context.Context, o domain.Organization) (domain.Organization, error) {
		// Mapear la entidad del dominio a la tabla
		tableOrg := mapper.MapOrganizationToTable(ctx, o.Name)

		var account table.Account
		err := conn.Where("email = ?", o.Operator.Contact.PrimaryEmail).Find(&account).Error
		if err != nil {
			return domain.Organization{}, err
		}
		var accountOrg table.AccountOrganization
		err = conn.Transaction(func(tx *gorm.DB) error {
			if err := conn.DB.Create(&tableOrg).Error; err != nil {
				return errors.Wrap(ErrOrganizationDatabase, "failed to create organization")
			}
			accountOrg.AccountID = account.ID
			accountOrg.OrganizationID = tableOrg.ID
			if err := conn.DB.Create(&accountOrg).Error; err != nil {
				return errors.Wrap(ErrOrganizationDatabase, "failed to link account to organization")
			}
			return nil
		})
		if err != nil {
			return domain.Organization{}, err
		}
		// Mapear de vuelta a la entidad de dominio
		savedOrg := domain.Organization{
			ID:      tableOrg.ID,
			Country: o.Country,
			Name:    o.Name,
		}
		return savedOrg, nil
	}
}
