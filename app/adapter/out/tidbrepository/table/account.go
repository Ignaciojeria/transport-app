package table

import (
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID         int64  `gorm:"primaryKey"`
	Email      string `gorm:"type:varchar(255);not null;unique"`
	DocumentID string `gorm:"type:char(64);uniqueIndex"`
}

func (a Account) MapAccount() domain.Account {
	return domain.Account{
		Email: a.Email,
	}
}
