package table

import (
	"encoding/json"
	"fmt"
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONMap permite representar un JSON flexible como map[string]string
type JSONMap map[string]string

type Webhook struct {
	gorm.Model
	ID         int64     `gorm:"primaryKey"`
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID `gorm:"not null"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
	Payload    JSONMap   `gorm:"type:json"`
}

// Map convierte el Webhook desde la capa de persistencia a la entidad de dominio.
func (w Webhook) Map() (domain.Webhook, error) {
	var webhook domain.Webhook

	// Reconstruye el mapa original a partir del JSONMap
	reconstructedMap := make(map[string]interface{})
	for key, valueStr := range w.Payload {
		var value interface{}
		if err := json.Unmarshal([]byte(valueStr), &value); err != nil {
			return domain.Webhook{}, fmt.Errorf("failed to unmarshal field %s: %w", key, err)
		}
		reconstructedMap[key] = value
	}

	// Marshal el mapa reconstruido a bytes
	reconstructedBytes, err := json.Marshal(reconstructedMap)
	if err != nil {
		return domain.Webhook{}, fmt.Errorf("failed to marshal reconstructed webhook: %w", err)
	}

	// Unmarshal final al struct del dominio
	if err := json.Unmarshal(reconstructedBytes, &webhook); err != nil {
		return domain.Webhook{}, fmt.Errorf("failed to unmarshal to domain.Webhook: %w", err)
	}

	return webhook, nil
}