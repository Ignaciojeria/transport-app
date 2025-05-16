package table

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	ID     int64  `gorm:"primaryKey"`
	Status string `gorm:"not null"`
}
