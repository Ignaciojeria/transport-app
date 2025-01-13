package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
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

		// Iniciar la transacción
		err := conn.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Guardar la organización
			if err := tx.Create(&tableOrg).Error; err != nil {
				return ErrOrganizationDatabase.New(err.Error())
			}
			// Guardar la API key asociada
			apiKeyRecord := table.ApiKey{
				OrganizationID: tableOrg.ID,
				Key:            o.Key,
				Status:         "active",
			}
			if err := tx.Create(&apiKeyRecord).Error; err != nil {
				return ErrOrganizationDatabase.New(err.Error())
			}

			orgCountry := table.OrganizationCountry{
				OrganizationID: tableOrg.ID,
				Country:        o.Country.Alpha2(),
			}

			if err := tx.Create(&orgCountry).Error; err != nil {
				return ErrOrganizationDatabase.New(err.Error())
			}

			return nil
		})

		if err != nil {
			return domain.Organization{}, err
		}

		// Mapear de vuelta a la entidad de dominio
		savedOrg := domain.Organization{
			Country: o.Country,
			Name:    o.Name,
			Email:   o.Email,
			Key:     o.Key,
		}

		return savedOrg, nil
	}
}
