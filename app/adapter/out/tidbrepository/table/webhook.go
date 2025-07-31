package table

import (
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

func (w Webhook) Map() domain.Webhook {
	var webhook domain.Webhook
	// Aquí podrías deserializar el payload si es necesario
	return webhook
} 