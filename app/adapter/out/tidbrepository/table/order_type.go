package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderType struct {
	gorm.Model
	ID          int64     `gorm:"primaryKey"`
	DocumentID  string    `gorm:"type:char(64);uniqueIndex"`
	Type        string    `gorm:"type:varchar(191);not null;"`
	TenantID    uuid.UUID `gorm:"not null;"`
	Tenant      Tenant    `gorm:"foreignKey:TenantID"`
	Description string    `gorm:"type:text"`
}

func (o OrderType) Map() domain.OrderType {
	return domain.OrderType{
		Type:        o.Type,
		Description: o.Description,
	}
}
