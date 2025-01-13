package tidbrepository

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type EnsureOrganizationForCountry func(context.Context, domain.Organization) (domain.Organization, error)

func init() {
	ioc.Registry(
		NewEnsureOrganizationForCountry,
		tidb.NewTIDBConnection)
}

func NewEnsureOrganizationForCountry(conn tidb.TIDBConnection) EnsureOrganizationForCountry {
	return func(ctx context.Context, org domain.Organization) (domain.Organization, error) {
		// Validar si el API Key existe
		var apiKey table.ApiKey
		err := conn.WithContext(ctx).Where("`key` = ?", org.Key).First(&apiKey).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Si no existe, retornar un error
				return domain.Organization{}, fmt.Errorf("API key not found: %s", org.Key)
			}
			// Retornar cualquier otro error de la base de datos
			return domain.Organization{}, fmt.Errorf("error querying API key: %w", err)
		}

		// Validar si el API Key está disponible para el país
		var orgCountry table.OrganizationCountry
		err = conn.WithContext(ctx).
			Where("organization_id = ? AND country = ?", apiKey.OrganizationID, org.Country.Alpha2()).
			First(&orgCountry).Error
		if err == nil {
			// Si ya existe para el país, retornar nil (no es necesario crear uno nuevo)
			return domain.Organization{
				OrganizationCountryID: orgCountry.ID,
			}, nil
		} else if err != gorm.ErrRecordNotFound {
			// Si ocurre otro error, retornarlo
			return domain.Organization{}, fmt.Errorf("error querying organization country: %w", err)
		}

		// Crear un nuevo registro en la tabla OrganizationCountry
		newOrgCountry := table.OrganizationCountry{
			OrganizationID: apiKey.OrganizationID,
			Country:        org.Country.Alpha2(),
		}
		if err := conn.WithContext(ctx).Create(&newOrgCountry).Error; err != nil {
			return domain.Organization{}, fmt.Errorf("error creating organization country: %w", err)
		}
		org.OrganizationCountryID = orgCountry.ID
		return org, nil
	}
}
