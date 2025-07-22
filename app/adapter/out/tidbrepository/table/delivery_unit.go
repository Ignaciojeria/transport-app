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
	// Simplified fields (matching optimization domain)
	Volume          int64          `gorm:"type:bigint;default:0;"`
	WeightValue     int64          `gorm:"type:bigint;default:0;"`
	InsuranceValue  int64          `gorm:"type:bigint;default:0;"`
	// Legacy complex fields (for backward compatibility)
	JSONDimensions  JSONDimensions `gorm:"type:json"`
	JSONWeight      JSONWeight     `gorm:"type:json"`
	JSONInsurance   JSONInsurance  `gorm:"type:json"`
	JSONItems       JSONItems      `gorm:"type:json"`
}

func (p DeliveryUnit) Map() domain.DeliveryUnit {
	deliveryUnit := domain.DeliveryUnit{
		Lpn:        p.Lpn,
		Dimensions: p.JSONDimensions.Map(),
		Weight:     p.JSONWeight.Map(),
		Insurance:  p.JSONInsurance.Map(),
		Items:      p.JSONItems.Map(),
	}
	
	// Set simplified values
	deliveryUnit.SetSimpleValues(p.Volume, p.WeightValue, p.InsuranceValue)
	
	return deliveryUnit
}
