package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

func init() {
	ioc.Registry(
		NewSaveOrganization,
		tidb.NewTIDBConnection,
	)
}

type SaveOrganization func(
	context.Context,
	domain.Organization,
) (domain.Organization, error)

func NewSaveOrganization(conn tidb.TIDBConnection) SaveOrganization {
	return func(ctx context.Context, o domain.Organization) (domain.Organization, error) {
		// Mapear la entidad del dominio a la tabla
		tableOrg := mapper.MapOrganizationToTable(o)

		// Intentar guardar la organización en la base de datos
		if err := conn.DB.Create(&tableOrg).Error; err != nil {
			return domain.Organization{}, ErrOrganizationDatabase.New(err.Error())
		}

		// Mapear de vuelta a la entidad de dominio
		savedOrg := domain.Organization{
			Country: o.Country,
			Name:    o.Name,
			Email:   o.Email,
			Key:     o.Key, // Si el key viene de otro lugar, podrías usar el generado en la DB si es necesario
		}

		return savedOrg, nil
	}
}
