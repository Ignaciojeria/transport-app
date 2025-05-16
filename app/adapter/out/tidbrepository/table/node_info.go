package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NodeInfo struct {
	gorm.Model
	ID          int64     `gorm:"primaryKey"`
	DocumentID  string    `gorm:"type:char(64);uniqueIndex"`
	ReferenceID string    `gorm:"type:varchar(191);not null;"`
	TenantID    uuid.UUID `gorm:"not null;"`
	Tenant      Tenant    `gorm:"foreignKey:TenantID"`
	Name        string    `gorm:"type:varchar(191);"`

	// Store document hashes without enforcing constraints
	NodeTypeDoc string   `gorm:"type:char(64)"`
	NodeType    NodeType `gorm:"-"` // Use "-" to tell GORM to ignore this field for DB operations

	ContactDoc string  `gorm:"type:char(64)"`
	Contact    Contact `gorm:"-"` // Ignore relationship for DB operations

	AddressInfoDoc string      `gorm:"type:char(64)"`
	AddressInfo    AddressInfo `gorm:"-"` // Ignore relationship for DB operations

	AddressLine2   string
	NodeReferences JSONReference `gorm:"type:json"`
}

func (n NodeInfo) Map() domain.NodeInfo {
	nodeInfo := domain.NodeInfo{
		ReferenceID:  domain.ReferenceID(n.ReferenceID),
		Name:         n.Name,
		References:   n.NodeReferences.Map(),
		AddressLine2: n.AddressLine2,
	}
	return nodeInfo
}
