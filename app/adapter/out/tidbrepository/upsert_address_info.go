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

type UpsertAddressInfo func(context.Context, domain.AddressInfo) (domain.AddressInfo, error)

func init() {
	ioc.Registry(NewUpsertAddressInfo, tidb.NewTIDBConnection)
}
func NewUpsertAddressInfo(conn tidb.TIDBConnection) UpsertAddressInfo {
	return func(ctx context.Context, ai domain.AddressInfo) (domain.AddressInfo, error) {
		var addressInfo table.AddressInfo
		err := conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("raw_address = ? AND organization_id = ?", ai.RawAddress(), ai.Organization.ID).
			First(&addressInfo).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.AddressInfo{}, err
		}
		addrInfoUpdated := addressInfo.Map().UpdateIfChanged(ai)
		dbAddressInfoToUpsert := mapper.MapAddressInfoTable(addrInfoUpdated, ai.Organization.ID)
		dbAddressInfoToUpsert.CreatedAt = addressInfo.CreatedAt
		if err := conn.Save(&dbAddressInfoToUpsert).Error; err != nil {
			return domain.AddressInfo{}, err
		}
		return dbAddressInfoToUpsert.Map(), nil
	}
}
