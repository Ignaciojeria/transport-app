package table

import (
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID       int64  `gorm:"primaryKey"`
	Email    string `gorm:"type:varchar(255);not null;unique"`
	IsActive bool   `gorm:"default:null"`
}

func (a Account) MapOperator() domain.Operator {
	return domain.Operator{}
}
