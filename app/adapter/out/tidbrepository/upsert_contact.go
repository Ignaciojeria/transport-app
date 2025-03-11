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

type UpsertContact func(ctx context.Context, c domain.Contact) (domain.Contact, error)

func init() {
	ioc.Registry(NewUpsertContact, tidb.NewTIDBConnection)
}

func NewUpsertContact(conn tidb.TIDBConnection) UpsertContact {
	return func(ctx context.Context, c domain.Contact) (domain.Contact, error) {
		var contact table.Contact
		err := conn.DB.WithContext(ctx).
			Table("contacts").
			Where("full_name = ? AND email = ? AND phone = ? AND national_id = ? AND organization_id = ?",
				c.FullName, c.Email, c.Phone, c.NationalID, c.Organization.ID).
			First(&contact).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Contact{}, err
		}
		contactWithChanges := contact.Map().UpdateIfChanged(c)
		dbContactToUpsert := mapper.MapContactToTable(contactWithChanges, c.Organization.ID)
		dbContactToUpsert.CreatedAt = contact.CreatedAt
		if err := conn.Omit("Organization").
			Save(&dbContactToUpsert).Error; err != nil {
			return domain.Contact{}, err
		}
		return dbContactToUpsert.Map(), nil
	}
}
