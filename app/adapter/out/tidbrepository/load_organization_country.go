package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type LoadOrganizationCountry func(context.Context, domain.Organization) (domain.Organization, error)

func init() {
	ioc.Registry(NewLoadOrganizationCountry, tidb.NewTIDBConnection)
}

func NewLoadOrganizationCountry(conn tidb.TIDBConnection) LoadOrganizationCountry {
	return func(ctx context.Context, o domain.Organization) (domain.Organization, error) {
		//var org table.Organization
		/*
			err := conn.DB.WithContext(ctx).
				Joins("JOIN organization_countries ON organization_countries.organization_id = organizations.id").
				Where("organization_countries.id = ?", o.OrganizationCountryID).
				First(&org).Error

			if err != nil {
				return domain.Organization{}, err
			}

			var orgCountry table.OrganizationCountry
			err = conn.DB.WithContext(ctx).
				Where("id = ?", o.OrganizationCountryID).
				First(&orgCountry).Error

			if err != nil {
				return domain.Organization{}, err
			}

			// Obtener la API key activa
			var apiKey table.ApiKey
			err = conn.DB.WithContext(ctx).
				Where("organization_id = ? AND status = ?", org.ID, "active").
				First(&apiKey).Error

			if err != nil {
				return domain.Organization{}, err
			}

			return domain.Organization{
				ID:                    org.ID,
				Name:                  org.Name,
				Email:                 org.Email,
				OrganizationCountryID: orgCountry.ID,
				Country:               orgCountry.Map().Country,
				Key:                   apiKey.Key,
			}, nil*/
		return domain.Organization{}, nil
	}
}
