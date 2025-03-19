package usecase

import (
	"context"
	"transport-app/app/adapter/out/firebaseauth"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type Register func(ctx context.Context, input domain.UserCredentials) error

func init() {
	ioc.Registry(NewRegister, firebaseauth.NewRegister, tidbrepository.NewUpsertAccount)
}
func NewRegister(register firebaseauth.Register, upsertAccount tidbrepository.UpsertAccount) Register {
	return func(ctx context.Context, input domain.UserCredentials) error {
		err := upsertAccount(ctx, domain.Account{
			Email: input.Email,
		})
		if err != nil {
			return err
		}
		return register(ctx, input)
	}
}
