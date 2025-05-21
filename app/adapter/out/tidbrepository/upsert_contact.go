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

type UpsertContact func(ctx context.Context, c domain.Contact) error

func init() {
	ioc.Registry(NewUpsertContact, database.NewConnectionFactory)
}

func NewUpsertContact(conn database.ConnectionFactory) UpsertContact {
	return func(ctx context.Context, c domain.Contact) error {
		var existing table.Contact
		err := conn.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", c.DocID(ctx)).
			First(&existing).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No existe → insert
			newContact := mapper.MapContactToTable(ctx, c)
			return conn.Omit("Tenant").Create(&newContact).Error
		}

		// Ya existe → update solo si cambió algo
		updated, changed := existing.Map().UpdateIfChanged(c)
		if !changed {
			return nil
		}

		updateData := mapper.MapContactToTable(ctx, updated)
		updateData.ID = existing.ID // necesario para que GORM haga UPDATE
		updateData.CreatedAt = existing.CreatedAt

		return conn.Omit("Tenant").Save(&updateData).Error
	}
}
