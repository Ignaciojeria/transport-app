package table

import "gorm.io/gorm"

type Province struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"default:null"`
	CountryCode string `gorm:"not null"`
	DocumentID  string `gorm:"type:char(64);uniqueIndex"`
}
