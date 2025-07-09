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

type SaveTenant func(ctx context.Context, tenant domain.Tenant) (domain.Tenant, error)

func NewSaveTenant(conn database.ConnectionFactory) SaveTenant {
	return func(ctx context.Context, tenant domain.Tenant) (domain.Tenant, error) {
		var existing table.Tenant
		err := conn.DB.WithContext(ctx).
			Table("tenants").
			Where("id = ?", tenant.ID).
			First(&existing).Error

		if err == nil {
			return domain.Tenant{
				ID:      existing.ID,
				Name:    existing.Name,
				Country: tenant.Country,
			}, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Tenant{}, err
		}

		tableTenant := mapper.MapTenantTable(ctx, tenant)
		tableTenant.ID = tenant.ID

		if err := conn.DB.Create(&tableTenant).Error; err != nil {
			return domain.Tenant{}, errors.Wrap(ErrTenantDatabase, "failed to create tenant")
		}

		return domain.Tenant{
			ID:      tableTenant.ID,
			Name:    tableTenant.Name,
			Country: tenant.Country,
		}, nil
	}
}
