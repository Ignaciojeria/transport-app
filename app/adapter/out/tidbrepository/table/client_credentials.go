package table

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ClientCredential struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	TenantID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	Tenant        Tenant         `gorm:"foreignKey:TenantID"`
	ClientID      string         `gorm:"uniqueIndex;not null"`
	ClientSecret  string         `gorm:"not null"`
	AllowedScopes pq.StringArray `gorm:"type:text[]"`
	Status        string         `gorm:"default:active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	ExpiresAt     *time.Time     `gorm:"type:timestamp"`
}
