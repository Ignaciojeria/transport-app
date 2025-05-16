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
		NewSaveTenant,
		database.NewConnectionFactory,
	)
}

type SaveTenant func(
	context.Context,
	domain.Tenant,
) (domain.Tenant, error)

func NewSaveTenant(conn database.ConnectionFactory) SaveTenant {
	return func(ctx context.Context, o domain.Tenant) (domain.Tenant, error) {
		// Mapear la entidad del dominio a la tabla
		tableOrg := mapper.MapTenantTable(ctx, o.Name)
		tableOrg.ID = o.ID
		var account table.Account
		err := conn.Where("email = ?", o.Operator.Contact.PrimaryEmail).Find(&account).Error
		if err != nil {
			return domain.Tenant{}, err
		}
		var accountOrg table.AccountTenant
		err = conn.Transaction(func(tx *gorm.DB) error {
			if err := conn.DB.Create(&tableOrg).Error; err != nil {
				return errors.Wrap(ErrTenantDatabase, "failed to create organization")
			}
			accountOrg.AccountID = account.ID
			accountOrg.TenantID = tableOrg.ID
			if err := conn.DB.Create(&accountOrg).Error; err != nil {
				return errors.Wrap(ErrTenantDatabase, "failed to link account to organization")
			}
			return nil
		})
		if err != nil {
			return domain.Tenant{}, err
		}
		// Mapear de vuelta a la entidad de dominio
		savedOrg := domain.Tenant{
			ID:      tableOrg.ID,
			Country: o.Country,
			Name:    o.Name,
		}
		return savedOrg, nil
	}
}
