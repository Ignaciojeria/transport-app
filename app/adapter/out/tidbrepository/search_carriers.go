package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"
	"transport-app/app/shared/sharedcontext"

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
			Joins("JOIN organizations org ON carriers.organization_id = org.id"). // Se une directamente con organizations
			Where("org.id = ?", sharedcontext.TenantIDFromContext(ctx)).          // Filtra solo por organization_id
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
