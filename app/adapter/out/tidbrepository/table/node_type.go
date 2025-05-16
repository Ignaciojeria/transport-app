package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NodeType struct {
	gorm.Model
	ID         int64     `gorm:"primaryKey"`
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID `gorm:"not null;"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
	Value      string    `gorm:"type:varchar(191);"`
}

func (n NodeType) Map() domain.NodeType {
	return domain.NodeType{
		Value: n.Value,
	}
}
