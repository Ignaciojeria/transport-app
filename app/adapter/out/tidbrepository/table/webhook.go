package table

import (
	"encoding/json"
	"log"
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Webhook struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:char(36);primaryKey"`
	Payload string    `gorm:"type:text"`
}

// Map converts the table Webhook to domain.Webhook by deserializing the Payload field
func (w Webhook) Map() domain.Webhook {
	var domainWebhook domain.Webhook
	
	if w.Payload != "" {
		if err := json.Unmarshal([]byte(w.Payload), &domainWebhook); err != nil {
			log.Printf("Error unmarshaling webhook payload: %v", err)
			// Return empty struct on unmarshal error, but log the error
			return domain.Webhook{}
		}
	}
	
	// Set the ID from the table field
	domainWebhook.ID = w.ID
	
	return domainWebhook
}