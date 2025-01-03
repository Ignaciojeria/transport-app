package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateOrganizationKey func(ctx context.Context, org domain.Organization) (domain.Organization, error)

func init() {
	ioc.Registry(
		NewCreateOrganizationKey,
		tidbrepository.NewSaveOrganization)
}

func NewCreateOrganizationKey(saveOrg tidbrepository.SaveOrganization) CreateOrganizationKey {
	return func(ctx context.Context, org domain.Organization) (domain.Organization, error) {
		return saveOrg(ctx, org)
	}
}
