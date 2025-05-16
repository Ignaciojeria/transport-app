package table

import (
	"encoding/json"
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Modelo de Vehículo
type Vehicle struct {
	gorm.Model
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID `gorm:"uniqueIndex:idx_vehicle_ref_org;uniqueIndex:idx_vehicle_plate_org"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
	Plate      string    `gorm:"type:varchar(20);uniqueIndex:idx_vehicle_plate_org"`
	IsActive   bool

	CertificateDate string `gorm:"type:varchar(50)"`

	VehicleCategoryDoc string          `gorm:"type:char(64);index"`
	VehicleCategory    VehicleCategory `gorm:"-"`

	VehicleHeadersDoc string         `gorm:"type:char(64);not null"`
	VehicleHeaders    VehicleHeaders `gorm:"-"`

	CarrierDoc string  `gorm:"type:char(64);index"`
	Carrier    Carrier `gorm:"-"`

	Weight          JSONB `gorm:"type:json"` // Tipo JSON para serializar Weight
	Insurance       JSONB `gorm:"type:json"` // Tipo JSON para serializar Insurance
	TechnicalReview JSONB `gorm:"type:json"` // Tipo JSON para serializar TechnicalReview
	Dimensions      JSONB `gorm:"type:json"` // Tipo JSON para serializar Dimensions
}

func (v Vehicle) Map() domain.Vehicle {
	vehicle := domain.Vehicle{
		Plate:           v.Plate,
		CertificateDate: v.CertificateDate,
	}

	// Deserializar los campos JSON
	if v.Weight != nil {
		_ = json.Unmarshal(v.Weight, &vehicle.Weight)
	}
	if v.Insurance != nil {
		_ = json.Unmarshal(v.Insurance, &vehicle.Insurance)
	}
	if v.TechnicalReview != nil {
		_ = json.Unmarshal(v.TechnicalReview, &vehicle.TechnicalReview)
	}
	if v.Dimensions != nil {
		_ = json.Unmarshal(v.Dimensions, &vehicle.Dimensions)
	}

	return vehicle
}
