package table

import (
	"encoding/json"
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type EventOutbox struct {
	gorm.Model
	ID             int64  `gorm:"primaryKey"`
	ReferenceID    string `gorm:"not null;"`
	EventType      string `gorm:"not null;"`
	EntityType     string `gorm:"not null;"`
	OrganizationID int64  `gorm:"not null;"`
	Attributes     []byte `gorm:"type:json"`
	Payload        []byte `gorm:"type:json"`
	Status         string `gorm:"default:'pending'"` // Valores posibles: pending, failed, processed
}

func MapEventOutbox(outbox domain.Outbox) EventOutbox {
	attrsBytes, _ := json.Marshal(outbox.Attributes)
	return EventOutbox{
		ReferenceID:    outbox.Attributes["referenceID"],
		OrganizationID: outbox.Organization.ID,
		Payload:        outbox.Payload,
		Status:         outbox.Status,
		EventType:      outbox.Attributes["eventType"],
		EntityType:     outbox.Attributes["entityType"],
		Attributes:     attrsBytes,
	}
}
