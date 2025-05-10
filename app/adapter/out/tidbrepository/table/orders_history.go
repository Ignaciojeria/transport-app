package table

import (
	"gorm.io/gorm"
)

type OrdersHistory struct {
	gorm.Model
	ID                   int64 `gorm:"primaryKey"`
	Channel              string
	OrderDoc             string             `gorm:"type:char(64);index"`
	OrderStatusID        int64              `gorm:"not null;index"`
	OrderStatus          Status             `gorm:"foreignKey:OrderStatusID"`
	RouteDoc             string             `gorm:"type:char(64);index"`
	VehicleDoc           string             `gorm:"type:char(64);index"`
	CarrierDoc           string             `gorm:"type:char(64);index"`
	DriverDoc            string             `gorm:"type:char(64);index"`
	NonDeliveryReasonDoc string             `gorm:"type:char(64);index"`
	EvidencePhotos       JSONEvidencePhotos `gorm:"type:json;default:'{}'"`
	RecipientFullName    string             `gorm:"default:''"`
	RecipientNationalID  string             `gorm:"default:''"`
	ExpectedQuantity     int
	DeliveredQuantity    int
}
