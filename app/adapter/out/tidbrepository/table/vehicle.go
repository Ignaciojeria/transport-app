package table

import (
	"transport-app/app/domain"

	"gorm.io/gorm"
)

// Modelo de Vehículo
type Vehicle struct {
	gorm.Model
	ID                int64  `gorm:"primaryKey"`
	TenantID          string `gorm:"not null;index;uniqueIndex:idx_vehicle_ref_org;uniqueIndex:idx_vehicle_plate_org"`
	Tenant            Tenant `gorm:"foreignKey:TenantID"`
	ReferenceID       string `gorm:"type:varchar(50);uniqueIndex:idx_vehicle_ref_org"`
	Plate             string `gorm:"type:varchar(20);uniqueIndex:idx_vehicle_plate_org"`
	IsActive          bool
	CertificateDate   string
	VehicleCategoryID *int64          `gorm:"default null;index"`
	VehicleCategory   VehicleCategory `gorm:"foreignKey:VehicleCategoryID"`
	VehicleHeadersID  int64           `gorm:"not null"`
	VehicleHeaders    VehicleHeaders  `gorm:"foreignKey:VehicleHeadersID"`
	Weight            JSONB           `gorm:"type:json"`          // Tipo JSON para serializar Weight
	Insurance         JSONB           `gorm:"type:json"`          // Tipo JSON para serializar Insurance
	TechnicalReview   JSONB           `gorm:"type:json"`          // Tipo JSON para serializar TechnicalReview
	Dimensions        JSONB           `gorm:"type:json"`          // Tipo JSON para serializar Dimensions
	CarrierID         *int64          `gorm:"default null;index"` // Relación con Carrier
	Carrier           Carrier         `gorm:"foreignKey:CarrierID"`
}

func (v Vehicle) Map() domain.Vehicle {
	return domain.Vehicle{
		//ReferenceID:     v.ReferenceID,
		Plate:           v.Plate,
		CertificateDate: v.CertificateDate,
		VehicleCategory: domain.VehicleCategory{
			Type:                v.VehicleCategory.Type,
			MaxPackagesQuantity: v.VehicleCategory.MaxPackagesQuantity,
		},
		Headers: v.VehicleHeaders.Map(),
		Carrier: v.Carrier.Map(),
	}
}
