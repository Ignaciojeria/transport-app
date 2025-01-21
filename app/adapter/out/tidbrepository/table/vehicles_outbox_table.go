package table

import (
	"encoding/json"
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type VehiclesOutbox struct {
	gorm.Model
	ID                    int64  `gorm:"primaryKey"`
	ReferenceID           string `gorm:"not null;index:idx_orders_outbox_unique,unique"`
	EventType             string `gorm:"not null;index:idx_orders_outbox_unique,unique"`
	OrganizationCountryID int64  `gorm:"not null;index:idx_orders_outbox_unique,unique"`
	Attributes            []byte `gorm:"type:json"`
	Payload               []byte `gorm:"type:json"`
	Status                string `gorm:"default:'pending'"` // Valores posibles: pending, failed, processed
}

func MapVehiclesOutbox(outbox domain.Outbox) VehiclesOutbox {
	attrsBytes, _ := json.Marshal(outbox.Attributes)
	return VehiclesOutbox{
		ReferenceID:           outbox.Attributes["referenceID"],
		OrganizationCountryID: outbox.Organization.OrganizationCountryID,
		Payload:               outbox.Payload,
		Status:                outbox.Status,
		EventType:             outbox.Attributes["eventType"],
		Attributes:            attrsBytes,
	}
}
