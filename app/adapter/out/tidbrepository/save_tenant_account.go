package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type SaveTenantAccount func(ctx context.Context, tenantAccount domain.TenantAccount, fsmState ...domain.FSMState) error

func init() {
	ioc.Registry(
		NewSaveTenantAccount,
		database.NewConnectionFactory,
		NewSaveFSMTransition,
	)
}

func NewSaveTenantAccount(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) SaveTenantAccount {
	return func(ctx context.Context, tenantAccount domain.TenantAccount, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var dbAccount table.Account
			err := tx.
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
				Invited:   tenantAccount.Invited,
				JoinedAt:  tenantAccount.JoinedAt,
			}

			// Upsert: actualizar si existe, crear si no existe
			result := tx.WithContext(ctx).
				Where("account_id = ? AND tenant_id = ?", dbAccount.ID, tenantAccount.Tenant.ID).
				Assign(link).
				FirstOrCreate(&link)

			if result.Error != nil {
				return errors.Wrap(ErrTenantDatabase, "failed to upsert account-tenant link")
			}

			// Persistir FSMState si está presente
			if len(fsmState) > 0 {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})
	}
}
