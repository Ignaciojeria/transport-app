package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"gorm.io/gorm"
)

type AddressInfo struct {
	gorm.Model
	ID           int64     `gorm:"primaryKey"`
	TenantID     uuid.UUID `gorm:"not null;"`
	Tenant       Tenant    `gorm:"foreignKey:TenantID"`
	DocumentID   string    `gorm:"type:char(64);uniqueIndex"`
	StateDoc     string    `gorm:"type:char(64);default:null"`
	State        State     `gorm:"-"`
	ProvinceDoc  string    `gorm:"type:char(64);default:null"`
	Province     Province  `gorm:"-"`
	DistrictDoc  string    `gorm:"type:char(64);default:null"`
	District     District  `gorm:"-"`
	AddressLine1 string    `gorm:"not null"`
	Latitude     float64   `gorm:"default:null"`
	Longitude    float64   `gorm:"default:null"`
	ZipCode      string    `gorm:"default:null"`
	TimeZone     string    `gorm:"default:null"`
}

func (a AddressInfo) Map() domain.AddressInfo {
	return domain.AddressInfo{
		State:        domain.State(a.State.Name),
		Province:     domain.Province(a.Province.Name),
		District:     domain.District(a.District.Name),
		AddressLine1: a.AddressLine1,
		Location:     orb.Point{a.Longitude, a.Latitude},
		ZipCode:      a.ZipCode,
		TimeZone:     a.TimeZone,
	}
}
