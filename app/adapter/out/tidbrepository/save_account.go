package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(NewSaveAccount, tidb.NewTIDBConnection)
}

type SaveAccount func(
	context.Context,
	domain.Account) (domain.Account, error)

func NewSaveAccount(conn tidb.TIDBConnection) SaveAccount {
	return func(ctx context.Context, o domain.Account) (domain.Account, error) {
		return domain.Account{}, nil
	}
}
