package table

import (
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type Route struct {
	gorm.Model
	ID                 int64                `gorm:"primaryKey;autoIncrement"`
	ReferenceID        string               `gorm:"type:varchar(255);not null"`
	TenantID           int64                `gorm:"default:null"`
	Tenant             Tenant               `gorm:"foreignKey:TenantID"`
	EndNodeReferenceID string               `gorm:"default:null"`
	JSONEndLocation    JSONRouteEndLocation `gorm:"type:json;default:null"`
	PlanID             int64                `gorm:"default:null"`
	Plan               Plan                 `gorm:"foreignKey:PlanID"`
	VehicleID          *int64               `gorm:"default:null"`
	Vehicle            Vehicle              `gorm:"foreignKey:VehicleID"`
	Driver             Driver
	CarrierID          *int64  `gorm:"default:null"`
	Carrier            Carrier `gorm:"foreignKey:CarrierID"`
}

func (r Route) Map() domain.Route {
	return domain.Route{
		Vehicle: r.Vehicle.Map(),
	}
}
