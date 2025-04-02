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

type UpsertContact func(ctx context.Context, c domain.Contact) error

func init() {
	ioc.Registry(NewUpsertContact, tidb.NewTIDBConnection)
}

func NewUpsertContact(conn tidb.TIDBConnection) UpsertContact {
	return func(ctx context.Context, c domain.Contact) error {
		var existing table.Contact
		err := conn.DB.WithContext(ctx).
			Table("contacts").
			Preload("Organization").
			Where("reference_id = ?", c.DocID()).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newContact := mapper.MapContactToTable(c, c.Organization.ID)
			return conn.Omit("Organization").Create(&newContact).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(c)
		if !changed {
			return nil
		}

		updateData := mapper.MapContactToTable(updated, c.Organization.ID)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Organization").Save(&updateData).Error
	}
}
