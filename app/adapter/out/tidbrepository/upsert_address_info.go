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
			// Upsert State
			if err := upsertState(ctx, tx, ai.State); err != nil {
				return err
			}

			// Upsert Province
			if err := upsertProvince(ctx, tx, ai.Province); err != nil {
				return err
			}

			// Upsert District
			if err := upsertDistrict(ctx, tx, ai.District); err != nil {
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

func upsertState(ctx context.Context, tx *gorm.DB, state domain.State) error {
	var existing table.State
	err := tx.WithContext(ctx).
		Table("states").
		Where("document_id = ?", state.DocID(ctx).String()).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Si ya existe, no hacemos nada porque el nombre es el mismo
	if err == nil {
		return nil
	}

	// No existe → insert
	newState := mapper.MapStateTable(ctx, state)
	return tx.Create(&newState).Error
}

func upsertProvince(ctx context.Context, tx *gorm.DB, province domain.Province) error {
	var existing table.Province
	err := tx.WithContext(ctx).
		Table("provinces").
		Where("document_id = ?", province.DocID(ctx).String()).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Si ya existe, no hacemos nada porque el nombre es el mismo
	if err == nil {
		return nil
	}

	// No existe → insert
	newProvince := mapper.MapProvinceTable(ctx, province)
	return tx.Create(&newProvince).Error
}

func upsertDistrict(ctx context.Context, tx *gorm.DB, district domain.District) error {
	var existing table.District
	err := tx.WithContext(ctx).
		Table("districts").
		Where("document_id = ?", district.DocID(ctx).String()).
		First(&existing).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Si ya existe, no hacemos nada porque el nombre es el mismo
	if err == nil {
		return nil
	}

	// No existe → insert
	newDistrict := mapper.MapDistrictTable(ctx, district)
	return tx.Create(&newDistrict).Error
}
