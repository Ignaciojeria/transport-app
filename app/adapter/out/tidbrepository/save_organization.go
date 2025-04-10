package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"
	"transport-app/app/shared/sharedcontext"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
	"gorm.io/gorm"
)

func init() {
	ioc.Registry(
		NewSaveOrganization,
		tidb.NewTIDBConnection,
	)
}

type SaveOrganization func(
	context.Context,
	domain.Operator,
	string,
) (domain.Organization, error)

func NewSaveOrganization(conn tidb.TIDBConnection) SaveOrganization {
	return func(ctx context.Context, o domain.Operator, orgName string) (domain.Organization, error) {
		// Mapear la entidad del dominio a la tabla
		tableOrg := mapper.MapOrganizationToTable(ctx, orgName)

		var account table.Account
		err := conn.Where("email = ?", o.Contact.PrimaryEmail).Find(&account).Error
		if err != nil {
			return domain.Organization{}, err
		}
		var accountOrg table.AccountOrganization
		err = conn.Transaction(func(tx *gorm.DB) error {
			if err := conn.DB.Create(&tableOrg).Error; err != nil {
				return ErrOrganizationDatabase.New("failed to create organization: %v", err)
			}
			accountOrg.AccountID = account.ID
			accountOrg.OrganizationID = tableOrg.ID
			if err := conn.DB.Create(&accountOrg).Error; err != nil {
				return ErrOrganizationDatabase.New("failed to link account to organization: %v", err)
			}
			return nil
		})
		if err != nil {
			return domain.Organization{}, err
		}
		// Mapear de vuelta a la entidad de dominio
		savedOrg := domain.Organization{
			ID:      tableOrg.ID,
			Country: countries.ByName(sharedcontext.TenantCountryFromContext(ctx)),
			Name:    orgName,
		}
		return savedOrg, nil
	}
}
