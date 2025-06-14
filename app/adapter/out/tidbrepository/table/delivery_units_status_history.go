package table

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeliveryUnitsStatusHistory struct {
	gorm.Model
	ID                           int64              `gorm:"primaryKey"`
	TenantID                     uuid.UUID          `gorm:"not null"`
	Tenant                       Tenant             `gorm:"foreignKey:TenantID"`
	DocumentID                   string             `gorm:"type:char(64);"`
	Channel                      string             `gorm:"default:''"`
	OrderDoc                     string             `gorm:"type:char(64);index"`
	DeliveryUnitDoc              string             `gorm:"type:char(64);index"`
	DeliveryUnitStatusDoc        string             `gorm:"type:char(64);index"`
	RouteDoc                     string             `gorm:"type:char(64);index"`
	NonDeliveryReasonReferenceID string             `json:"non_delivery_reason_reference_id"`
	NonDeliveryReason            string             `json:"non_delivery_reason"`
	NonDeliveryDetail            string             `json:"non_delivery_detail"`
	EvidencePhotos               JSONEvidencePhotos `gorm:"type:json;"`
	RecipientFullName            string             `gorm:"default:''"`
	RecipientNationalID          string             `gorm:"default:''"`
	ConfirmDeliveryHandledAt     time.Time          `gorm:"default:null"`
	ConfirmDeliveryLatitude      float64            `gorm:"default:0"`
	ConfirmDeliveryLongitude     float64            `gorm:"default:0"`
	ManualChangePerformedBy      string             `gorm:"default:''"`
	ManualChangeReason           string             `gorm:"default:''"`
}
