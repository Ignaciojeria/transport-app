package table

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	ID         int64  `gorm:"primaryKey"`
	DocumentID string `gorm:"type:char(64);uniqueIndex"`
	Status     string `gorm:"not null"`
}
