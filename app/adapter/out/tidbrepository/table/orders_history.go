package table

import (
	"gorm.io/gorm"
)

type OrdersHistory struct {
	gorm.Model
	ID                   int64              `gorm:"primaryKey"`
	OrderDoc             string             `gorm:"type:char(64);index"`
	OrderStatusID        int64              `gorm:"not null;index"`
	OrderStatus          Status             `gorm:"foreignKey:OrderStatusID"`
	RouteDoc             string             `gorm:"type:char(64);index"`
	VehicleDoc           string             `gorm:"type:char(64);index"`
	CarrierDoc           string             `gorm:"type:char(64);index"`
	DriverDoc            string             `gorm:"type:char(64);index"`
	NonDeliveryReasonDoc string             `gorm:"type:char(64);index"`
	OriginNodeDoc        string             `gorm:"type:char(64);index"`
	EvidencePhotos       JSONEvidencePhotos `gorm:"type:json;default:'{}'"`
	RecipientFullName    string             `gorm:"default:''"`
	RecipientNationalID  string             `gorm:"default:''"`

	/*

		OrderReferenceID          int64              `gorm:"default:0"`
		OrderStatusID             int64              `gorm:"default:0"`
		OrderStatus               OrderStatus        `gorm:"foreignKey:OrderStatusID"`
		Address                   string             `gorm:"default:''"`
		ContainerLPN              string             `gorm:"default:''"`
		PackageLPN                string             `gorm:"default:''"`
		PlanConsumer              string             `gorm:"default:''"`
		PlanCommerce              string             `gorm:"default:''"`
		OrderCommerce             string             `gorm:"default:''"`
		OrderConsumer             string             `gorm:"default:''"`
		Province                  string             `gorm:"default:''"`
		District                  string             `gorm:"default:''"`
		State                     string             `gorm:"default:''"`
		Volume                    float64            `gorm:"default:0"`
		WorkingHoursStart         time.Time          `gorm:"default:'2000-01-01 00:00:00'"`
		WorkingHoursEnd           time.Time          `gorm:"default:'2000-01-01 00:00:00'"`
		TravelType                string             `gorm:"default:''"`
		OrderType                 string             `gorm:"default:''"`
		GeocodingSource           string             `gorm:"default:''"`
		VehiclePlate              string             `gorm:"default:''"`
		CarrierNationalID         string             `gorm:"default:''"`
		CarrierName               string             `gorm:"default:''"`
		DriverEmail               string             `gorm:"default:''"`
		DriverName                string             `gorm:"default:''"`
		DriverRut                 string             `gorm:"default:''"`
		RouteReferenceID          string             `gorm:"default:''"`
		OrderDeliverySequence     string             `gorm:"default:''"`
		NodeName                  string             `gorm:"default:''"`
		CheckoutRejectionID       int64              `gorm:"default:0"`
		CheckoutRejection         CheckoutRejection  `gorm:"foreignKey:CheckoutRejectionID"`
		CheckoutRejectionComment  string             `gorm:"default:''"`
		EvidencePhotos            JSONEvidencePhotos `gorm:"type:json;default:'{}'"`
		LegDistance               float64            `gorm:"default:0"`
		GeocodedLatitude          float64            `gorm:"default:0"`
		GeocodedLongitude         float64            `gorm:"default:0"`
		DriverDeliveryLatitude    float64            `gorm:"default:0"`
		DriverDeliveryLongitude   float64            `gorm:"default:0"`
		RecipientFullName         string             `gorm:"default:''"`
		RecipientNationalID       string             `gorm:"default:''"`
		DeliveryMethod            string             `gorm:"default:''"`
		PromisedDate              time.Time          `gorm:"default:'2000-01-01 00:00:00'"`
		DeliveryAttemptDateTime   time.Time          `gorm:"default:'2000-01-01 00:00:00'"`
		EstimatedDeliveryDateTime time.Time          `gorm:"default:'2000-01-01 00:00:00'"`
		RouteStartedDateTime      time.Time          `gorm:"default:'2000-01-01 00:00:00'"`
	*/
}
