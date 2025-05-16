package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Driver struct {
	gorm.Model
	DocumentID string    `gorm:"type:char(64);uniqueIndex"`
	TenantID   uuid.UUID `gorm:"uniqueIndex:idx_driver_ref_org;uniqueIndex:idx_driver_national_org"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID"`
	Name       string    `gorm:"not null"`
	NationalID string    `gorm:"type:varchar(20);default:null;uniqueIndex:idx_driver_national_org"`
	Email      string    `gorm:"type:varchar(191);default:null"`
}

func (d Driver) Map() domain.Driver {
	return domain.Driver{
		Name:       d.Name,
		NationalID: d.NationalID,
		Email:      d.Email,
	}
}
