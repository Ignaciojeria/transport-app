package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Province struct {
	gorm.Model
	ID         int64     `gorm:"primaryKey"`
	Name       string    `gorm:"type:varchar(191);not null"`
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID `gorm:"not null"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
}

func (p Province) Map() domain.Province {
	return domain.Province(p.Name)
}
