package table

import (
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type State struct {
	gorm.Model
	ID         int64  `gorm:"primaryKey"`
	Name       string `gorm:"type:varchar(191);not null"`
	DocumentID string `gorm:"type:char(64);uniqueIndex"`
	TenantID   string `gorm:"not null"`
}

func (s State) Map() domain.State {
	return domain.State(s.Name)
}
