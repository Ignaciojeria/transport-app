package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/pkg/errors"
)

type SaveTenantAccount func(ctx context.Context, tenantAccount domain.TenantAccount) error

func init() {
	ioc.Registry(
		NewSaveTenantAccount,
		database.NewConnectionFactory,
	)
}

func NewSaveTenantAccount(conn database.ConnectionFactory) SaveTenantAccount {
	return func(ctx context.Context, tenantAccount domain.TenantAccount) error {
		var dbAccount table.Account
		err := conn.
			DB.
			WithContext(ctx).
			Where("email = ?", tenantAccount.Account.Email).
			First(&dbAccount).Error
		if err != nil {
			return errors.Wrap(err, "account not found")
		}

		link := table.AccountTenant{
			AccountID: dbAccount.ID,
			Role:      tenantAccount.Role,
			Status:    tenantAccount.Status,
			TenantID:  tenantAccount.Tenant.ID,
		}

		if err := conn.DB.WithContext(ctx).Create(&link).Error; err != nil {
			return errors.Wrap(ErrTenantDatabase, "failed to link account to tenant")
		}

		return nil
	}
}
