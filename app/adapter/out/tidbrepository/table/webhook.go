package table

import (
	"encoding/json"
	"fmt"
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Webhook struct {
	gorm.Model
	ID         int64  `gorm:"primaryKey"`
	DocumentID string `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID `gorm:"not null;"`
	Tenant     Tenant `gorm:"foreignKey:TenantID"`
	Payload    JSONMap `gorm:"type:json"`
}

func (w Webhook) Map() (domain.Webhook, error) {
	var webhook domain.Webhook
	
	// Reconstruct the original JSON from the JSONMap
	reconstructedMap := make(map[string]interface{})
	for key, valueStr := range w.Payload {
		var value interface{}
		if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
			return domain.Webhook{}, fmt.Errorf("failed to unmarshal field %s: %w", key, err)
		}
		reconstructedMap[key] = value
	}
	
	// Marshal the reconstructed map back to JSON
	reconstructedBytes, err := json.Marshal(reconstructedMap)
	if err != nil {
		return domain.Webhook{}, fmt.Errorf("failed to marshal reconstructed webhook: %w", err)
	}
	
	// Unmarshal into the domain.Webhook struct
	if err := json.Unmarshal(reconstructedBytes, &webhook); err != nil {
		return domain.Webhook{}, fmt.Errorf("failed to unmarshal to domain.Webhook: %w", err)
	}
	
	return webhook, nil
} 