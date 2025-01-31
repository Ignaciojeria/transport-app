package tidbrepository

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

type UpsertOrderHeaders func(context.Context, domain.Headers) (domain.Headers, error)

func init() {
	ioc.Registry(NewUpsertOrderHeaders, tidb.NewTIDBConnection)
}
func NewUpsertOrderHeaders(conn tidb.TIDBConnection) UpsertOrderHeaders {
	return func(ctx context.Context, h domain.Headers) (domain.Headers, error) {
		var orderHeaders table.OrderHeaders
		err := conn.DB.WithContext(ctx).
			Table("order_headers").
			Where(
				"commerce = ? AND consumer = ? AND organization_country_id = ?",
				h.Commerce, h.Consumer, h.Organization.OrganizationCountryID).
			First(&orderHeaders).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Headers{}, err
		}
		if orderHeaders.Commerce != "" {
			return orderHeaders.Map(), nil
		}
		orderHeadersTbl := mapper.MapOrderHeaders(h)
		if err := conn.Save(&orderHeadersTbl).Error; err != nil {
			return domain.Headers{}, err
		}
		return orderHeadersTbl.Map(), nil
	}
}
