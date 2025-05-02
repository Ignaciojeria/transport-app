package table

import (
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type NonDeliveryReason struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey"`
	DocumentID  string `gorm:"type:char(64);uniqueIndex"`
	ReferenceID string
	Reason      string
}

func (t NonDeliveryReason) Map() domain.NonDeliveryReason {
	return domain.NonDeliveryReason{
		ReferenceID: t.ReferenceID,
		Reason:      t.Reason,
	}
}
