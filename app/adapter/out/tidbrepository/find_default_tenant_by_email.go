package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type FindDefaultTenantByEmail func(context.Context, string) (domain.TenantAccount, error)

func init() {
	ioc.Registry(NewFindDefaultTenantByEmail, database.NewConnectionFactory)
}

func NewFindDefaultTenantByEmail(conn database.ConnectionFactory) FindDefaultTenantByEmail {
	return func(ctx context.Context, email string) (domain.TenantAccount, error) {
		var accountTenant table.AccountTenant
		err := conn.DB.WithContext(ctx).
			Joins("JOIN accounts ON account_tenants.account_id = accounts.id").
			Joins("JOIN tenants ON account_tenants.tenant_id = tenants.id").
			Preload("Tenant").
			Where("accounts.email = ? AND account_tenants.invited = ?", email, false).
			First(&accountTenant).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.TenantAccount{}, nil
			}
			return domain.TenantAccount{}, err
		}

		return domain.TenantAccount{
			Tenant: accountTenant.Tenant.Map(),
			Account: domain.Account{
				Email: email,
			},
			Role:     accountTenant.Role,
			Status:   accountTenant.Status,
			Invited:  accountTenant.Invited,
			JoinedAt: accountTenant.JoinedAt,
		}, nil
	}
}
