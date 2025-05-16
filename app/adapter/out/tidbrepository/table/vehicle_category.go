package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleCategory struct {
	gorm.Model
	ID                  int64     `gorm:"primaryKey"`
	DocumentID          string    `gorm:"type:char(64);uniqueIndex"`
	TenantID            uuid.UUID `gorm:"not null;uniqueIndex:idx_org_country_type"`
	Tenant              Tenant    `gorm:"foreignKey:TenantID"`
	Type                string    `gorm:"type:varchar(191);not null;uniqueIndex:idx_org_country_type"`
	MaxPackagesQuantity int
}

func (vc VehicleCategory) Map() domain.VehicleCategory {
	return domain.VehicleCategory{
		MaxPackagesQuantity: vc.MaxPackagesQuantity,
		Type:                vc.Type,
	}
}
