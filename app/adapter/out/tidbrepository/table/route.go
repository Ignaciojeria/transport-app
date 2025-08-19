package table

import (
	"encoding/json"
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Route struct {
	gorm.Model
	ID                int64                `gorm:"primaryKey"`
	Raw               json.RawMessage      `gorm:"type:jsonb;default:null"`
	ReferenceID       string               `gorm:"type:varchar(255);not null"`
	DocumentID        string               `gorm:"type:char(64);uniqueIndex"`
	TenantID          uuid.UUID            `gorm:"not null"`
	Tenant            Tenant               `gorm:"foreignKey:TenantID"`
	EndNodeInfoDoc    string               `gorm:"default:null"`
	EndNodeInfo       NodeInfo             `gorm:"-"`
	OriginNodeInfoDoc string               `gorm:"type:char(64);index"`
	OriginNodeInfo    NodeInfo             `gorm:"-"`
	JSONEndLocation   JSONRouteEndLocation `gorm:"type:json;default:null"`
	PlanDoc           string               `gorm:"type:char(64);index"`
	Plan              Plan                 `gorm:"-"`
	VehicleDoc        string               `gorm:"type:char(64);index"`
	Vehicle           Vehicle              `gorm:"-"`
	DriverDoc         string               `gorm:"type:char(64);index"`
	Driver            Driver               `gorm:"-"`
	CarrierDoc        string               `gorm:"type:char(64);index"`
	Carrier           Carrier              `gorm:"-"`
}

func (r Route) Map() domain.Route {
	return domain.Route{
		ReferenceID: r.ReferenceID,
		Origin:      r.OriginNodeInfo.Map(),
		Vehicle:     r.Vehicle.Map(),
	}
}
