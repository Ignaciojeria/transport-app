package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Carrier struct {
	gorm.Model
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID `gorm:"uniqueIndex:idx_carrier_ref_org;uniqueIndex:idx_carrier_national_org"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
	Name       string    `gorm:"not null"`
	NationalID string    `gorm:"type:varchar(20);default:null;uniqueIndex:idx_carrier_national_org"`
}

func (c Carrier) Map() domain.Carrier {
	return domain.Carrier{
		Name:       c.Name,
		NationalID: c.NationalID,
	}
}
