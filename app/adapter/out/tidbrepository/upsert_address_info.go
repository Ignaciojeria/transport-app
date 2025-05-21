package tidbrepository

import (
	"context"
	"errors"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/adapter/out/tidbrepository/table/mapper"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type UpsertAddressInfo func(context.Context, domain.AddressInfo) error

func init() {
	ioc.Registry(NewUpsertAddressInfo, database.NewConnectionFactory)
}

func NewUpsertAddressInfo(conn database.ConnectionFactory) UpsertAddressInfo {
	return func(ctx context.Context, ai domain.AddressInfo) error {
		var existing table.AddressInfo
		err := conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", ai.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newAddressInfo := mapper.MapAddressInfoTable(ctx, ai)
			return conn.Omit("Tenant").Create(&newAddressInfo).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(ai)
		if !changed {
			return nil // No hay cambios, no hacemos nada
		}

		updateData := mapper.MapAddressInfoTable(ctx, updated)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Tenant").Save(&updateData).Error
	}
}
