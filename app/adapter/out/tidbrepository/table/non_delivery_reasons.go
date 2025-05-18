package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NonDeliveryReason struct {
	gorm.Model
	ID          int64     `gorm:"primaryKey"`
	DocumentID  string    `gorm:"type:char(64);uniqueIndex"`
	TenantID    uuid.UUID `gorm:"not null;index"`
	Tenant      Tenant    `gorm:"foreignKey:TenantID"`
	ReferenceID string
	Reason      string
	Details     string
}

func (t NonDeliveryReason) Map() domain.NonDeliveryReason {
	return domain.NonDeliveryReason{
		ReferenceID: t.ReferenceID,
		Reason:      t.Reason,
		Details:     t.Details,
	}
}
