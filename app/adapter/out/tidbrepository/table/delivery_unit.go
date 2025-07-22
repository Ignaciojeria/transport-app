package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryUnit struct {
	gorm.Model
	SizeCategoryDoc string         `gorm:"type:char(64);not null;"`
	SizeCategory    SizeCategory   `gorm:"-"`
	ID              int64          `gorm:"primaryKey"`
	TenantID        uuid.UUID      `gorm:"not null;"`
	Tenant          Tenant         `gorm:"foreignKey:TenantID"`
	DocumentID      string         `gorm:"type:char(64);uniqueIndex"`
	Lpn             string         `gorm:"type:varchar(191);not null;"`
	Volume          int64          `gorm:"type:bigint;default:0;"`
	JSONDimensions  JSONDimensions `gorm:"type:json"`
	JSONWeight      JSONWeight     `gorm:"type:json"`
	JSONInsurance   JSONInsurance  `gorm:"type:json"`
	JSONItems       JSONItems      `gorm:"type:json"`
}

func (p DeliveryUnit) Map() domain.DeliveryUnit {
	return domain.DeliveryUnit{
		Lpn:        p.Lpn,
		Volume:     p.Volume,
		Dimensions: p.JSONDimensions.Map(),
		Weight:     p.JSONWeight.Map(),
		Insurance:  p.JSONInsurance.Map(),
		Items:      p.JSONItems.Map(),
	}
}
