package table

import (
	"encoding/json"
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type NodesOutbox struct {
	gorm.Model
	ID                    int64  `gorm:"primaryKey"`
	ReferenceID           string `gorm:"not null;"`
	EventType             string `gorm:"not null;"`
	OrganizationCountryID int64  `gorm:"not null;"`
	Attributes            []byte `gorm:"type:json"`
	Payload               []byte `gorm:"type:json"`
	Status                string `gorm:"default:'pending'"` // Valores posibles: pending, failed, processed
}

func MapNodesOutbox(outbox domain.Outbox) NodesOutbox {
	attrsBytes, _ := json.Marshal(outbox.Attributes)
	return NodesOutbox{
		ReferenceID:           outbox.Attributes["referenceID"],
		OrganizationCountryID: outbox.Organization.OrganizationCountryID,
		Payload:               outbox.Payload,
		Status:                outbox.Status,
		EventType:             outbox.Attributes["eventType"],
		Attributes:            attrsBytes,
	}
}
