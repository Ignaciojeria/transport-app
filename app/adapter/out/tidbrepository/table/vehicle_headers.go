package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleHeaders struct {
	gorm.Model
	ID         int64     `gorm:"primaryKey"`
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	Commerce   string    `gorm:"default:null"`
	Consumer   string    `gorm:"default:null"`
	Channel    string    `gorm:"default:null"`
	TenantID   uuid.UUID `gorm:"not null;index;"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
}

func (m VehicleHeaders) Map() domain.Headers {
	return domain.Headers{
		Consumer: m.Consumer,
		Commerce: m.Commerce,
		Channel:  m.Channel,
	}
}
