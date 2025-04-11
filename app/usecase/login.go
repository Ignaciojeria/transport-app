package usecase

import (
	"context"
	"transport-app/app/adapter/out/firebaseauth"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type Login func(ctx context.Context, userCreds domain.UserCredentials) (domain.ProviderToken, error)

func init() {
	ioc.Registry(
		NewLogin,
		tidbrepository.NewFindOrganizationByEmail,
		firebaseauth.NewLogin)
}

func NewLogin(
	findOrganizationByEmail tidbrepository.FindOrganizationByEmail,
	login firebaseauth.Login,
) Login {
	return func(ctx context.Context, userCreds domain.UserCredentials) (domain.ProviderToken, error) {
		return login(ctx, userCreds)
	}
}
