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

type UpsertAddressInfo func(context.Context, domain.AddressInfo, ...domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertAddressInfo, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertAddressInfo(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertAddressInfo {
	return func(ctx context.Context, ai domain.AddressInfo, fsmState ...domain.FSMState) error {
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
				if err := tx.Omit("Tenant").Create(&newAddressInfo).Error; err != nil {
					return err
				}
			} else {
				// Ya existe → update solo si cambió algo
				updated, changed := existing.Map().UpdateIfChanged(ai)
				if changed {
					updateData := mapper.MapAddressInfoTable(ctx, updated)
					updateData.ID = existing.ID // necesario para que GORM haga UPDATE
					updateData.CreatedAt = existing.CreatedAt
					updateData.PoliticalAreaDoc = ai.PoliticalArea.DocID(ctx).String()
					if err := tx.Omit("Tenant").Save(&updateData).Error; err != nil {
						return err
					}
				}
			}

			// Guardar FSMState si se provee
			if len(fsmState) > 0 && saveFSMTransition != nil {
				return saveFSMTransition(ctx, fsmState[0], tx)
			}
			return nil
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
