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

type UpsertCarrier func(context.Context, domain.Carrier) (domain.Carrier, error)

func init() {
	ioc.Registry(NewUpsertCarrier, tidb.NewTIDBConnection)
}
func NewUpsertCarrier(conn tidb.TIDBConnection) UpsertCarrier {
	return func(ctx context.Context, c domain.Carrier) (domain.Carrier, error) {
		carrier := table.Carrier{}
		err := conn.DB.WithContext(ctx).Table("carriers").
			Where("reference_id = ? AND organization_id = ?",
				string("TODO"), c.Organization.ID).First(&carrier).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Carrier{}, err
		}
		carrierWithChanges, _ := carrier.Map().UpdateIfChanged(c)
		DBCarrierToUpsert := mapper.MapCarrierToTable(carrierWithChanges)
		DBCarrierToUpsert.CreatedAt = carrier.CreatedAt
		if err := conn.Save(&DBCarrierToUpsert).Error; err != nil {
			return domain.Carrier{}, err
		}
		return DBCarrierToUpsert.Map(), nil
	}
}
