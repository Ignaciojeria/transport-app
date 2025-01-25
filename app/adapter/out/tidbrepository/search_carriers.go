package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type SearchCarriers func(context.Context, domain.CarrierSearchFilters) ([]domain.Carrier, error)

func init() {
	ioc.Registry(
		NewSearchCarriers,
		tidb.NewTIDBConnection)
}
func NewSearchCarriers(conn tidb.TIDBConnection) SearchCarriers {
	return func(ctx context.Context, csf domain.CarrierSearchFilters) ([]domain.Carrier, error) {
		var carriers []table.Carrier
		if err := conn.DB.WithContext(ctx).
			Joins("JOIN organization_countries oc ON carriers.organization_country_id = oc.id").
			Joins("JOIN organizations org ON oc.organization_id = org.id").
			Joins("JOIN api_keys ak ON org.id = ak.organization_id").
			Where("ak.key = ? AND oc.country = ?", csf.Organization.Key, csf.Organization.Country.Alpha2()).
			Limit(csf.Size).
			Offset(csf.Page).
			Find(&carriers).Error; err != nil {
			return nil, err
		}

		result := make([]domain.Carrier, 0, len(carriers))
		for _, v := range carriers {
			result = append(result, v.Map())
		}
		return result, nil
	}
}
