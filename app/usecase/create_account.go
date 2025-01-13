package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateAccount func(ctx context.Context, account domain.Account) (domain.Account, error)

func init() {
	ioc.Registry(
		NewCreateAccount,
		tidbrepository.NewEnsureOrganizationForCountry,
		tidbrepository.NewSaveAccount)
}

func NewCreateAccount(
	ensureOrganizationForCountry tidbrepository.EnsureOrganizationForCountry,
	saveAccount tidbrepository.SaveAccount) CreateAccount {
	return func(ctx context.Context, e domain.Account) (domain.Account, error) {
		if _, err := ensureOrganizationForCountry(ctx, e.Organization); err != nil {
			return domain.Account{}, err
		}
		return saveAccount(ctx, e)
	}
}
