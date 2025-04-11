package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateAccount func(ctx context.Context, account domain.Operator) (domain.Operator, error)

func init() {
	ioc.Registry(
		NewCreateAccount,
		tidbrepository.NewEnsureOrganizationForCountry,
		tidbrepository.NewSaveAccount)
}

func NewCreateAccount(
	ensureOrganizationForCountry tidbrepository.EnsureOrganizationForCountry,
	saveAccount tidbrepository.SaveAccount) CreateAccount {
	return func(ctx context.Context, e domain.Operator) (domain.Operator, error) {
		return saveAccount(ctx, e)
	}
}
