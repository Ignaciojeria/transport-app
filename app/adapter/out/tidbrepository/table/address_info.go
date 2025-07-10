package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"gorm.io/gorm"
)

type AddressInfo struct {
	gorm.Model
	ID                      int64         `gorm:"primaryKey"`
	TenantID                uuid.UUID     `gorm:"not null;"`
	Tenant                  Tenant        `gorm:"foreignKey:TenantID"`
	DocumentID              string        `gorm:"type:char(64);uniqueIndex"`
	PoliticalAreaDoc        string        `gorm:"type:char(64);default:null"`
	PoliticalArea           PoliticalArea `gorm:"-"`
	AddressLine1            string        `gorm:"not null"`
	AddressLine2            string        `gorm:"default:null"`
	Latitude                float64       `gorm:"default:null"`
	Longitude               float64       `gorm:"default:null"`
	CoordinateSource        string        `gorm:"default:null"`
	CoordinateConfidence    float64       `gorm:"default:null"`
	CoordinateMessage       string        `gorm:"default:null"`
	CoordinateReason        string        `gorm:"default:null"`
	PoliticalAreaConfidence float64       `gorm:"default:null"`
	PoliticalAreaMessage    string        `gorm:"default:null"`
	PoliticalAreaReason     string        `gorm:"default:null"`
	ZipCode                 string        `gorm:"default:null"`
}

func (a AddressInfo) Map() domain.AddressInfo {
	return domain.AddressInfo{
		PoliticalArea: domain.PoliticalArea{
			Code:            a.PoliticalArea.Code,
			AdminAreaLevel1: a.PoliticalArea.AdminAreaLevel1,
			AdminAreaLevel2: a.PoliticalArea.AdminAreaLevel2,
			AdminAreaLevel3: a.PoliticalArea.AdminAreaLevel3,
			AdminAreaLevel4: a.PoliticalArea.AdminAreaLevel4,
			TimeZone:        a.PoliticalArea.TimeZone,
		},
		AddressLine1: a.AddressLine1,
		AddressLine2: a.AddressLine2,
		Coordinates: domain.Coordinates{
			Point:  orb.Point{a.Longitude, a.Latitude},
			Source: a.CoordinateSource,
			Confidence: domain.CoordinatesConfidence{
				Level:   a.CoordinateConfidence,
				Message: a.CoordinateMessage,
				Reason:  a.CoordinateReason,
			},
		},
		ZipCode: a.ZipCode,
	}
}
