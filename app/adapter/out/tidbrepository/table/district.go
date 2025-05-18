package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type District struct {
	gorm.Model
	ID         int64     `gorm:"primaryKey"`
	Name       string    `gorm:"type:varchar(191);not null"`
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID `gorm:"not null"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
}

func (d District) Map() domain.District {
	return domain.District(d.Name)
}
