package table

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryUnitsStatusHistory struct {
	gorm.Model
	ID                       int64              `gorm:"primaryKey"`
	TenantID                 uuid.UUID          `gorm:"not null"`
	Tenant                   Tenant             `gorm:"foreignKey:TenantID"`
	DocumentID               string             `gorm:"type:char(64);"`
	Channel                  string             `gorm:"default:''"`
	OrderDoc                 string             `gorm:"type:char(64);index"`
	DeliveryUnitDoc          string             `gorm:"type:char(64);index"`
	DeliveryUnitStatusDoc    string             `gorm:"type:char(64);index"`
	PlanDoc                  string             `gorm:"type:char(64);index"`
	RouteDoc                 string             `gorm:"type:char(64);index"`
	VehicleDoc               string             `gorm:"type:char(64);index"`
	CarrierDoc               string             `gorm:"type:char(64);index"`
	DriverDoc                string             `gorm:"type:char(64);index"`
	NonDeliveryReasonDoc     string             `gorm:"type:char(64);index"`
	EvidencePhotos           JSONEvidencePhotos `gorm:"type:json;"`
	RecipientFullName        string             `gorm:"default:''"`
	RecipientNationalID      string             `gorm:"default:''"`
	ConfirmDeliveryHandledAt time.Time          `gorm:"default:null"`
	ConfirmDeliveryLatitude  float64            `gorm:"default:0"`
	ConfirmDeliveryLongitude float64            `gorm:"default:0"`
}
