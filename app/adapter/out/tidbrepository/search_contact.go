package tidbrepository

import (
	"context"
	"errors"
	"fmt"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/tidb"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"gorm.io/gorm"
)

type SearchContact func(ctx context.Context, c domain.Contact) (domain.Contact, error)

func init() {
	ioc.Registry(NewSearchContact, tidb.NewTIDBConnection)
}

func NewSearchContact(conn tidb.TIDBConnection) SearchContact {
	return func(ctx context.Context, c domain.Contact) (domain.Contact, error) {
		var contact table.Contact

		// Búsqueda del contacto
		err := conn.DB.WithContext(ctx).
			Table("contacts").
			Where("full_name = ? AND email = ? AND phone = ? AND organization_country_id = ?",
				c.FullName, c.Email, c.Phone, c.Organization.OrganizationCountryID).
			First(&contact).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.Contact{}, fmt.Errorf("contact not found: %w", err)
			}
			return domain.Contact{}, fmt.Errorf("error querying contact: %w", err)
		}

		return contact.Map(), nil
	}
}
