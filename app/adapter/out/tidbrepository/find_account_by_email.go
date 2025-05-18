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

type FindAccountByEmail func(context.Context, string) (domain.Account, error)

func init() {
	ioc.Registry(NewFindAccountByEmail, database.NewConnectionFactory)
}

func NewFindAccountByEmail(conn database.ConnectionFactory) FindAccountByEmail {
	return func(ctx context.Context, email string) (domain.Account, error) {
		var accountTbl table.Account
		err := conn.DB.WithContext(ctx).
			Table("accounts").
			Where("email = ?", email).
			First(&accountTbl).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.Account{}, nil
			}
			return domain.Account{}, err
		}

		return domain.Account{
			Email: accountTbl.Email,
		}, nil
	}
}
