package table

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey;autoIncrement"`
	StartNodeReferenceID  string              `gorm:"default:null"`
	JSONStartLocation     JSONPlanLocation    `gorm:"type:json;default:null"`
	OrganizationCountryID int64               `gorm:"not null"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	ReferenceID           string              `gorm:"type:varchar(255);not null"`
	PlannedDate           time.Time           `gorm:"type:date;not null"`
	PlanTypeID            int64               `gorm:"not null"`
	PlanType              PlanType            `gorm:"foreignKey:PlanTypeID"`
	PlanningStatusID      int64               `gorm:"not null"`
	PlanningStatus        PlanningStatus      `gorm:"foreignKey:PlanningStatusID"`
}

type PlanLocation struct {
	Longitude float64
	Latitude  float64
}

type JSONPlanLocation PlanLocation

// Implementamos los métodos necesarios para el manejo de JSON
func (j *JSONPlanLocation) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONPlanLocation value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// Value implementa la interfaz driver.Valuer para convertir la estructura en JSON
func (j JSONPlanLocation) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (p Plan) Map() domain.Plan {
	return domain.Plan{
		ID:             p.ID,
		ReferenceID:    p.ReferenceID,
		PlannedDate:    p.PlannedDate,
		PlanningStatus: p.PlanningStatus.Map(),
		PlanType:       p.PlanType.Map(),
	}
}

type PlanType struct {
	gorm.Model
	ID                    int64               `gorm:"type:bigint;primaryKey;autoIncrement"`
	OrganizationCountryID int64               `gorm:"type:bigint;not null;index;uniqueIndex:idx_plan_type_org_name"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID;references:ID"`
	Name                  string              `gorm:"type:varchar(100);not null;uniqueIndex:idx_plan_type_org_name"`
}

func (pt PlanType) Map() domain.PlanType {
	return domain.PlanType{
		Organization: pt.OrganizationCountry.Map(),
		ID:           pt.ID,
		Value:        pt.Name,
	}
}

type PlanningStatus struct {
	gorm.Model
	ID                    int64               `gorm:"type:bigint;primaryKey;autoIncrement"`
	OrganizationCountryID int64               `gorm:"type:bigint;not null;index;uniqueIndex:idx_planning_status_org_name"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID;references:ID"`
	Name                  string              `gorm:"type:varchar(100);not null;uniqueIndex:idx_planning_status_org_name"`
}

func (ps PlanningStatus) Map() domain.PlanningStatus {
	return domain.PlanningStatus{
		Organization: ps.OrganizationCountry.Map(),
		ID:           ps.ID,
		Value:        ps.Name,
	}
}

type Route struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey;autoIncrement"`
	ReferenceID           string              `gorm:"type:varchar(255);not null"`
	OrganizationCountryID int64               `gorm:"default:null"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	EndNodeReferenceID    string              `gorm:"default:null"`
	JSONEndLocation       JSONPlanLocation    `gorm:"type:json;default:null"`
	PlanID                int64               `gorm:"default:null"`
	Plan                  Plan                `gorm:"foreignKey:PlanID"`
	AccountID             *int64              `gorm:"default:null"`
	Account               Account             `gorm:"foreignKey:AccountID"`
	VehicleID             *int64              `gorm:"default:null"`
	Vehicle               Vehicle             `gorm:"foreignKey:VehicleID"`
	CarrierID             *int64              `gorm:"default:null"`
	Carrier               Carrier             `gorm:"foreignKey:CarrierID"`
}

func (r Route) Map() domain.Route {
	return domain.Route{
		Organization: r.OrganizationCountry.Map(),
		ReferenceID:  r.ReferenceID,
		ID:           r.ID,
		Vehicle:      r.Vehicle.Map(),
		Operator:     r.Account.MapOperator(),
		PlanID:       r.Plan.ID,
	}
}
