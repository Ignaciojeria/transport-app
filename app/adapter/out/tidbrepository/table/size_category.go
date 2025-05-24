package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SizeCategory struct {
	gorm.Model
	ID         int64     `gorm:"primaryKey"`
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID `gorm:"not null;index"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
	Code       string    `gorm:"type:varchar(64);index"`
}

func (s SizeCategory) Map() domain.SizeCategory {
	return domain.SizeCategory{
		Code: s.Code,
	}
}
