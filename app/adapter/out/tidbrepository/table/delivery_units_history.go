package table

import (
	"gorm.io/gorm"
)

type DeliveryUnitsHistory struct {
	gorm.Model
	ID                   int64              `gorm:"primaryKey"`
	DocumentID           string             `gorm:"type:char(64);uniqueIndex"`
	Channel              string             `gorm:"default:''"`
	OrderDoc             string             `gorm:"type:char(64);index"`
	DeliveryUnitDoc      string             `gorm:"type:char(64);index"`
	OrderStatusDoc       string             `gorm:"type:char(64);index"`
	PlanDoc              string             `gorm:"type:char(64);index"`
	RouteDoc             string             `gorm:"type:char(64);index"`
	VehicleDoc           string             `gorm:"type:char(64);index"`
	CarrierDoc           string             `gorm:"type:char(64);index"`
	DriverDoc            string             `gorm:"type:char(64);index"`
	NonDeliveryReasonDoc string             `gorm:"type:char(64);index"`
	EvidencePhotos       JSONEvidencePhotos `gorm:"type:json;default:'{}'"`
	RecipientFullName    string             `gorm:"default:''"`
	RecipientNationalID  string             `gorm:"default:''"`
}
