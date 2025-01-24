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

type UpsertConsumer func(context.Context, domain.Consumer) (domain.Consumer, error)

func init() {
	ioc.Registry(NewUpsertConsumer, tidb.NewTIDBConnection)
}
func NewUpsertConsumer(conn tidb.TIDBConnection) UpsertConsumer {
	return func(ctx context.Context, c domain.Consumer) (domain.Consumer, error) {
		var consumer table.Consumer
		err := conn.DB.WithContext(ctx).
			Table("consumers").
			Where("name = ? AND organization_country_id = ?", c.Value, c.Organization.OrganizationCountryID).
			First(&consumer).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Consumer{}, err
		}
		if consumer.Name != "" {
			return consumer.Map(), nil
		}
		consumer = mapper.MapConsumerToTable(consumer.Map().UpdateIfChanged(c))
		if err := conn.Save(&consumer).Error; err != nil {
			return domain.Consumer{}, err
		}
		return consumer.Map(), nil
	}
}
