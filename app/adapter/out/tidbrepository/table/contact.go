package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	ID         int64         `gorm:"primaryKey"`
	DocumentID string        `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID     `gorm:"not null;"`
	Tenant     Tenant        `gorm:"foreignKey:TenantID"`
	FullName   string        `gorm:"type:varchar(191);"`
	Email      string        `gorm:"type:varchar(191);"`
	Phone      string        `gorm:"type:varchar(191);"`
	NationalID string        `gorm:"type:varchar(191);"`
	Documents  JSONReference `gorm:"type:json"`
}

func (c Contact) Map() domain.Contact {
	return domain.Contact{
		FullName:     c.FullName,
		PrimaryEmail: c.Email,
		PrimaryPhone: c.Phone,
		NationalID:   c.NationalID,
		Documents:    c.Documents.MapDocuments(),
	}
}
