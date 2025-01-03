package usecase

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository"
	"transport-app/app/domain"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/google/uuid"
	"github.com/joomcode/errorx"
)

type CreateOrganizationKey func(ctx context.Context, org domain.Organization) (domain.Organization, error)

func init() {
	ioc.Registry(
		NewCreateOrganizationKey,
		tidbrepository.NewSaveOrganization,
		tidbrepository.NewFindOrganizationByEmail,
	)
}

func NewCreateOrganizationKey(
	saveOrg tidbrepository.SaveOrganization,
	findOrganizationByEmail tidbrepository.FindOrganizationByEmail,
) CreateOrganizationKey {
	return func(ctx context.Context, org domain.Organization) (domain.Organization, error) {
		// Buscar si la organización ya existe
		_, err := findOrganizationByEmail(ctx, org.Email)
		if err == nil {
			// Si la organización ya existe, retornar el error específico
			return domain.Organization{}, ErrOrganizationAlreadyExists.New("email already used")
		}

		// Si el error no es que no se encontró la organización, propagarlo
		if !errorx.IsOfType(err, tidbrepository.ErrOrganizationNotFound) {
			return domain.Organization{}, err
		}
		org.Key = uuid.NewString() + "-" + uuid.NewString()
		// Crear la organización si no existe
		return saveOrg(ctx, org)
	}
}
