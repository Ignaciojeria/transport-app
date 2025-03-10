package usecase

import (
	"context"
	"transport-app/app/adapter/out/firebaseauth"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type Register func(ctx context.Context, input domain.UserCredentials) error

func init() {
	ioc.Registry(NewRegister, firebaseauth.NewRegister)
}
func NewRegister(register firebaseauth.Register) Register {
	return func(ctx context.Context, input domain.UserCredentials) error {
		return register(ctx, input)
	}
}
