package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertAccount func(context.Context, domain.Account, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertAccount, database.NewConnectionFactory, NewSaveFSMTransition)
}
func NewUpsertAccount(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertAccount {
	return func(ctx context.Context, a domain.Account, fsmState ...domain.FSMState) error {
		return conn.Transaction(func(tx *gorm.DB) error {
			var accountTbl table.Account
			err := tx.WithContext(ctx).
				Table("accounts").
				Where("email = ?", a.Email).
				First(&accountTbl).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if err == nil {
				// La cuenta ya existe, solo persistir FSMState si está presente
				if len(fsmState) > 0 && saveFSMTransition != nil {
					return saveFSMTransition(ctx, fsmState[0], tx)
				}
				return nil
			}

			// La cuenta no existe, crearla
			accountTbl = mapper.MapAccountTable(a)
			if err := tx.Save(&accountTbl).Error; err != nil {
				return err
			}

			// Persistir FSMState si está presente
			if len(fsmState) > 0 && saveFSMTransition != nil {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}

			return nil
		})
	}
}
