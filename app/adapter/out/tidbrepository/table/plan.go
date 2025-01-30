package table

import (
	"transport-app/app/domain"

	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey;autoIncrement"`
	OrganizationCountryID int64               `gorm:"not null"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	ReferenceID           string              `gorm:"type:varchar(255);not null"`
	Date                  string              `gorm:"type:varchar(50);not null"`
	PlanTypeID            int64               `gorm:"not null"`
	PlanType              PlanType            `gorm:"foreignKey:PlanTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PlanningStatusID      int64               `gorm:"not null"`
	PlanningStatus        PlanningStatus      `gorm:"foreignKey:PlanningStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (p Plan) Map() domain.Plan {
	return domain.Plan{
		ID:             p.ID,
		ReferenceID:    p.ReferenceID,
		Date:           p.Date,
		PlanningStatus: p.PlanningStatus.Map(),
		PlanType:       p.PlanType.Map(),
	}
}

type PlanType struct {
	gorm.Model
	ID                    int64               `gorm:"primaryKey;autoIncrement"`
	OrganizationCountryID int64               `gorm:"not null"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Name                  string              `gorm:"type:varchar(100);not null;unique"`
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
	ID                    int64               `gorm:"primaryKey;autoIncrement"`
	OrganizationCountryID int64               `gorm:"not null"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	Name                  string              `gorm:"type:varchar(100);not null;unique"`
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
	OrganizationCountryID int64               `gorm:"not null"`
	OrganizationCountry   OrganizationCountry `gorm:"foreignKey:OrganizationCountryID"`
	PlanID                int                 `gorm:"not null"`
	Plan                  Plan                `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AccountID             int64               `gorm:"not null"`
	Account               Account             `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	VehicleID             int64               `gorm:"not null"`
	Vehicle               Vehicle             `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CarrierID             int64               `gorm:"not null"`
	Carrier               Carrier             `gorm:"foreignKey:CarrierID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (r Route) Map(orders []domain.Order) domain.Route {
	return domain.Route{
		Organization: r.OrganizationCountry.Map(),
		ID:           r.ID,
		Vehicle:      r.Vehicle.Map(),
		Operator:     r.Account.MapOperator(),
		Plan:         r.Plan.Map(),
		Orders:       orders,
	}
}
