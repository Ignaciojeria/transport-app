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

type UpsertAccount func(context.Context, domain.Operator) error

func init() {
	ioc.Registry(NewUpsertAccount, database.NewConnectionFactory)
}
func NewUpsertAccount(conn database.ConnectionFactory) UpsertAccount {
	return func(ctx context.Context, a domain.Operator) error {
		var accountTbl table.Account
		err := conn.DB.WithContext(ctx).
			Table("accounts").
			Where("email = ?", a.Contact.PrimaryEmail).
			First(&accountTbl).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err == nil {
			return nil
		}
		accountTbl = mapper.MapAccountTable(a)
		return conn.Save(&accountTbl).Error
	}
}
