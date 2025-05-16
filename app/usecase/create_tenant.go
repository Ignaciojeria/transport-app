package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type CreateTenant func(ctx context.Context, org domain.Tenant) (domain.Tenant, error)

func init() {
	ioc.Registry(
		NewCreateTenant,
		tidbrepository.NewSaveTenant,
	)
}

func NewCreateTenant(
	saveTenant tidbrepository.SaveTenant,
) CreateTenant {
	return func(ctx context.Context, org domain.Tenant) (domain.Tenant, error) {
		return saveTenant(ctx, org)
	}
}
