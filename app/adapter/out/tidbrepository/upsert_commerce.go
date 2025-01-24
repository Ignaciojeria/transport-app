package tidbrepository

/*
import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertCommerce func(context.Context, domain.Commerce) (domain.Commerce, error)

func init() {
	ioc.Registry(NewUpsertCommerce, tidb.NewTIDBConnection)
}
func NewUpsertCommerce(conn tidb.TIDBConnection) UpsertCommerce {
	return func(ctx context.Context, c domain.Commerce) (domain.Commerce, error) {
		var commerce table.Commerce
		err := conn.DB.WithContext(ctx).
			Table("commerces").
			Where("name = ? AND organization_country_id = ?", c.Value, c.Organization.OrganizationCountryID).
			First(&commerce).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Commerce{}, err
		}
		if commerce.Name != "" {
			return commerce.Map(), nil
		}
		commerce = mapper.MapCommerceToTable(commerce.Map().UpdateIfChanged(c))
		if err := conn.Save(&commerce).Error; err != nil {
			return domain.Commerce{}, err
		}
		return commerce.Map(), nil
	}
}
*/
