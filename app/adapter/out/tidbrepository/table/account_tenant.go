package table

import (
	"time"

	"github.com/google/uuid"
)

type AccountTenant struct {
	AccountID int64     `gorm:"primaryKey"`
	TenantID  uuid.UUID `gorm:"primaryKey"`
	Tenant    Tenant    `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
	Role      string    `gorm:"type:varchar(50);default:null"`
	Account   Account   `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
	Status    string
	Invited   bool      `gorm:"default:false"`
	JoinedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
