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
		return conn.DB.Transaction(func(tx *gorm.DB) error {
			// Upsert PoliticalArea
			if err := upsertPoliticalArea(ctx, tx, ai.PoliticalArea); err != nil {
				return err
			}

			var existing table.AddressInfo
			err := tx.WithContext(ctx).
				Table("address_infos").
				Where("document_id = ?", ai.DocID(ctx)).
				First(&existing).Error

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// No existe → insert
				newAddressInfo := mapper.MapAddressInfoTable(ctx, ai)
				return tx.Omit("Tenant").Create(&newAddressInfo).Error
			}

			// Ya existe → update solo si cambió algo
			updated, changed := existing.Map().UpdateIfChanged(ai)
			if !changed {
				return nil // No hay cambios, no hacemos nada
			}

			updateData := mapper.MapAddressInfoTable(ctx, updated)
			updateData.ID = existing.ID // necesario para que GORM haga UPDATE
			updateData.CreatedAt = existing.CreatedAt

			return tx.Omit("Tenant").Save(&updateData).Error
		})
	}
}

func upsertPoliticalArea(ctx context.Context, tx *gorm.DB, pa domain.PoliticalArea) error {
	var existing table.PoliticalArea
	err := tx.WithContext(ctx).
		Table("political_areas").
		Where("document_id = ?", pa.DocID(ctx).String()).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Si ya existe, no hacemos nada porque los datos son los mismos
	if err == nil {
		return nil
	}

	// No existe → insert
	newPoliticalArea := mapper.MapPoliticalArea(ctx, pa)
	return tx.Create(&newPoliticalArea).Error
}
