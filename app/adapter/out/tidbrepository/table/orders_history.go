package table

import "gorm.io/gorm"

type OrdersHistory struct {
	gorm.Model
	ID                  int64              `gorm:"primaryKey"`
	OrderReferenceID    int64              `gorm:"not null;index"`
	VehiclePlate        int64              `gorm:"not null;index"`
	CarrierNationalID   string             `gorm:"not null;index"`
	CarrierName         string             `gorm:"not null;index"`
	RouteReferenceID    string             `gorm:"not null;index"`
	OrderStatusID       int64              `gorm:"not null;index"`
	OrderStatus         OrderStatus        `gorm:"foreignKey:OrderStatusID"`
	CheckoutRejectionID int64              `gorm:"not null;index"`
	CheckoutRejection   CheckoutRejection  `gorm:"foreignKey:CheckoutRejectionID"`
	EvidencePhotos      JSONEvidencePhotos `gorm:"type:json"`
	NodeReferenceID     string
	NodeName            string
	Latitude            float64
	Longitude           float64
	RecipientFullName   string
	RecipientNationalID string
}
