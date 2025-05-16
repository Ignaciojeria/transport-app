package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryUnit struct {
	gorm.Model
	ID             int64          `gorm:"primaryKey"`
	TenantID       uuid.UUID      `gorm:"not null;"`
	Tenant         Tenant         `gorm:"foreignKey:TenantID"`
	DocumentID     string         `gorm:"type:char(64);uniqueIndex"`
	Lpn            string         `gorm:"type:varchar(191);not null;"`
	JSONDimensions JSONDimensions `gorm:"type:json"`
	JSONWeight     JSONWeight     `gorm:"type:json"`
	JSONInsurance  JSONInsurance  `gorm:"type:json"`
	JSONItems      JSONItems      `gorm:"type:json"`
}

func (p DeliveryUnit) Map() domain.Package {
	return domain.Package{
		Lpn:        p.Lpn,
		Dimensions: p.JSONDimensions.Map(),
		Weight:     p.JSONWeight.Map(),
		Insurance:  p.JSONInsurance.Map(),
		Items:      p.JSONItems.Map(),
	}
}
