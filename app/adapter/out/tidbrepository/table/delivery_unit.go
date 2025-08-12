package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryUnit struct {
	gorm.Model
	SizeCategoryDoc string       `gorm:"type:char(64);not null;"`
	SizeCategory    SizeCategory `gorm:"-"`
	ID              int64        `gorm:"primaryKey"`
	TenantID        uuid.UUID    `gorm:"not null;"`
	Tenant          Tenant       `gorm:"foreignKey:TenantID"`
	DocumentID      string       `gorm:"type:char(64);uniqueIndex"`
	Lpn             string       `gorm:"type:varchar(191);not null;"`
	// Simplified fields (matching optimization domain)
	Volume    int64     `gorm:"type:bigint;default:0;"`
	Weight    int64     `gorm:"type:bigint;default:0;"`
	Price     int64     `gorm:"type:bigint;default:0;"`
	JSONItems JSONItems `gorm:"type:json"`
}

func (p DeliveryUnit) Map() domain.DeliveryUnit {
	return domain.DeliveryUnit{
		Lpn:          p.Lpn,
		SizeCategory: p.SizeCategory.Map(),
		Volume:       &p.Volume,
		Weight:       &p.Weight,
		Price:        &p.Price,
		Items:        p.JSONItems.Map(),
	}
}
