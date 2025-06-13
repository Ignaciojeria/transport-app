package table

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Skill struct {
	gorm.Model
	ID         int64     `gorm:"primaryKey"`
	Name       string    `gorm:"type:varchar(255);not null"`
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID `gorm:"not null;index"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
}
