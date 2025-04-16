package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type EnsureOrganizationForCountry func(context.Context, domain.Organization) (domain.Organization, error)

func init() {
	ioc.Registry(
		NewEnsureOrganizationForCountry,
		database.NewConnectionFactory)
}
func NewEnsureOrganizationForCountry(conn database.ConnectionFactory) EnsureOrganizationForCountry {
	return func(ctx context.Context, org domain.Organization) (domain.Organization, error) {
		/*
					var result struct {
						OrganizationID        int64
						OrganizationName      string
						OrganizationCountryID int64
						Country               string
					}

					err := conn.WithContext(ctx).Raw(`
			            SELECT
			                o.id as organization_id,
			                o.name as organization_name,
			                oc.id as organization_country_id,
			                oc.country
			            FROM api_keys ak
			            JOIN organizations o ON o.id = ak.organization_id
			            JOIN organization_countries oc ON
			                oc.organization_id = ak.organization_id
			                AND oc.country = ?
			            WHERE ak.key = ?
			        `, org.Country.Alpha2(), org.Key).Scan(&result).Error

					if err != nil || result.OrganizationID == 0 {
						return domain.Organization{}, fmt.Errorf("organization not found for country %s and key %s", org.Country.Alpha2(), org.Key)
					}
					retrievedOrg := domain.Organization{
						ID:                    result.OrganizationID,
						Name:                  result.OrganizationName,
						OrganizationCountryID: result.OrganizationCountryID,
						Country:               countries.ByName(result.Country),
						Key:                   org.Key,
					}
					return retrievedOrg, nil*/
		return domain.Organization{}, nil
	}
}
