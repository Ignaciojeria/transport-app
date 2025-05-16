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
	State        string    `gorm:"default:null"`
	Province     string    `gorm:"default:null"`
	District     string    `gorm:"default:null"`
	AddressLine1 string    `gorm:"not null"`
	Latitude     float64   `gorm:"default:null"`
	Longitude    float64   `gorm:"default:null"`
	ZipCode      string    `gorm:"default:null"`
	TimeZone     string    `gorm:"default:null"`
}

func (a AddressInfo) Map() domain.AddressInfo {
	return domain.AddressInfo{
		State:    domain.State(a.State),
		Province: domain.Province(a.Province),
		//	Locality:     a.Locality,
		District:     domain.District(a.District),
		AddressLine1: a.AddressLine1,
		//	AddressLine2: a.AddressLine2,
		//	AddressLine3: a.AddressLine3,
		Location: orb.Point{a.Longitude, a.Latitude},
		ZipCode:  a.ZipCode,
		TimeZone: a.TimeZone,
	}
}
