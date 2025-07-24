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
		NewSaveFSMTransition,
	)
}

type SaveTenant func(ctx context.Context, tenant domain.Tenant, fsmState ...domain.FSMState) (domain.Tenant, error)

func NewSaveTenant(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) SaveTenant {
	return func(ctx context.Context, tenant domain.Tenant, fsmState ...domain.FSMState) (domain.Tenant, error) {
		var result domain.Tenant
		err := conn.Transaction(func(tx *gorm.DB) error {
			var existing table.Tenant
			err := tx.WithContext(ctx).
				Table("tenants").
				Where("id = ?", tenant.ID).
				First(&existing).Error

			if err == nil {
				// El tenant ya existe, solo persistir FSMState si está presente
				result = domain.Tenant{
					ID:      existing.ID,
					Name:    existing.Name,
					Country: tenant.Country,
				}
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			// El tenant no existe, crearlo
			tableTenant := mapper.MapTenantTable(ctx, tenant)
			tableTenant.ID = tenant.ID

			if err := tx.Create(&tableTenant).Error; err != nil {
				return errors.Wrap(ErrTenantDatabase, "failed to create tenant")
			}

			result = domain.Tenant{
				ID:      tableTenant.ID,
				Name:    tableTenant.Name,
				Country: tenant.Country,
			}

			// Persistir FSMState si está presente
			if len(fsmState) > 0 && saveFSMTransition != nil {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})

		return result, err
	}
}
