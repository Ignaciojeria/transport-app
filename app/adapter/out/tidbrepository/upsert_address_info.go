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

type UpsertAddressInfo func(context.Context, domain.AddressInfo) error

func init() {
	ioc.Registry(NewUpsertAddressInfo, tidb.NewTIDBConnection)
}

func NewUpsertAddressInfo(conn tidb.TIDBConnection) UpsertAddressInfo {
	return func(ctx context.Context, ai domain.AddressInfo) error {
		var existing table.AddressInfo
		err := conn.DB.WithContext(ctx).
			Table("address_infos").
			Preload("Organization").
			Where("reference_id = ?", ai.ReferenceID()).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newAddressInfo := mapper.MapAddressInfoTable(ai, ai.Organization.ID)
			return conn.Omit("Organization").Create(&newAddressInfo).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(ai)
		if !changed {
			return nil // No hay cambios, no hacemos nada
		}

		updateData := mapper.MapAddressInfoTable(updated, ai.Organization.ID)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Organization").Save(&updateData).Error
	}
}