package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateOrganization func(ctx context.Context, org domain.Operator) (domain.Organization, error)

func init() {
	ioc.Registry(
		NewCreateOrganization,
		tidbrepository.NewSaveOrganization,
	)
}

func NewCreateOrganization(
	saveOrg tidbrepository.SaveOrganization,
) CreateOrganization {
	return func(ctx context.Context, org domain.Operator) (domain.Organization, error) {
		return saveOrg(ctx, org)
	}
}
