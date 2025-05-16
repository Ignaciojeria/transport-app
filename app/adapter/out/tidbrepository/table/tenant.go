package table

import (
	"transport-app/app/domain"

	"github.com/biter777/countries"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name    string    `gorm:"type:varchar(255);not null;"`
	Country string    `gorm:"type:varchar(255);not null;"`
}

func (o Tenant) Map() domain.Tenant {
	return domain.Tenant{
		ID:      o.ID,
		Name:    o.Name,
		Country: countries.ByName(o.Country),
	}
}
