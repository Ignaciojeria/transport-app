package table

import (
	"time"
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	ReferenceID       string      `gorm:"type:varchar(255);not null"`
	DocumentID        string      `gorm:"type:char(64);uniqueIndex"`
	Name              string      `gorm:"default:null"`
	PlanHeadersDoc    string      `gorm:"type:char(64);not null"`
	PlanHeaders       PlanHeaders `gorm:"-"`
	OriginNodeInfoDoc string      `gorm:"type:char(64);index"`
	OriginNodeInfo    NodeInfo    `gorm:"-"`
	TenantID          uuid.UUID   `gorm:"not null"`
	Tenant            Tenant      `gorm:"foreignKey:TenantID"`
	PlannedDate       time.Time   `gorm:"type:date;not null"`
}

func (p Plan) Map() domain.Plan {
	return domain.Plan{
		ReferenceID: p.ReferenceID,
		PlannedDate: p.PlannedDate,
	}
}
