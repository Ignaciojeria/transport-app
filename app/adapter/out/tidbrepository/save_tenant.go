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
		// Buscar si ya existe
		var existing table.Tenant
		err := conn.DB.WithContext(ctx).
			Table("tenants").
			Where("id = ?", o.ID).
			First(&existing).Error

		// Si ya existe, retornar sin crear nuevamente
		if err == nil {
			return domain.Tenant{
				ID:      existing.ID,
				Name:    existing.Name,
				Country: o.Country, // Country puede no estar en la tabla, si es así ajusta
			}, nil
		}

		// Solo retornar si error distinto de "not found"
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Tenant{}, err
		}

		// Mapear y crear nuevo tenant
		tableOrg := mapper.MapTenantTable(ctx, o.Name)
		tableOrg.ID = o.ID

		var account table.Account
		err = conn.Where("email = ?", o.Operator.Contact.PrimaryEmail).Find(&account).Error
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

		return domain.Tenant{
			ID:      tableOrg.ID,
			Name:    tableOrg.Name,
			Country: o.Country,
		}, nil
	}
}
