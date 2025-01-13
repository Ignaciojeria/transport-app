package table

import (
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type OrdersOutbox struct {
	gorm.Model
	ID                    int64  `gorm:"primaryKey"`
	ReferenceID           string `gorm:"not null;index:idx_orders_outbox_unique,unique"`
	EventType             string `gorm:"not null;index:idx_orders_outbox_unique,unique"`
	OrganizationCountryID int64  `gorm:"not null;index:idx_orders_outbox_unique,unique"`
	Payload               []byte `gorm:"type:json"`
	Processed             bool   `gorm:"default:false"`
}

func MapOrderOutbox(outbox domain.Outbox) OrdersOutbox {
	return OrdersOutbox{
		ReferenceID:           outbox.ReferenceID,
		EventType:             outbox.EventType,
		OrganizationCountryID: outbox.Organization.ID, // Suponiendo que Organization tiene un campo ID.
		Payload:               outbox.Payload,
		Processed:             false,
	}
}
