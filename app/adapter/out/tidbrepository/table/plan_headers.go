package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PlanHeaders struct {
	gorm.Model
	ID         int64     `gorm:"primaryKey"`
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	Commerce   string    `gorm:"not null"`
	Consumer   string    `gorm:"not null"`
	Channel    string    `gorm:"not null"`
	TenantID   uuid.UUID `gorm:"not null;index;"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
}

func (m PlanHeaders) Map() domain.Headers {
	return domain.Headers{
		Consumer: m.Consumer,
		Commerce: m.Commerce,
		Channel:  m.Channel,
	}
}
