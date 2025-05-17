package table

import (
	"time"
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	ID                   int64     `gorm:"primaryKey;autoIncrement"`
	Name                 string    `gorm:"default:null"`
	StartNodeReferenceID string    `gorm:"default:null"`
	OriginNodeInfoDoc    string    `gorm:"type:char(64);index"`
	OriginNodeInfo       NodeInfo  `gorm:"-"`
	TenantID             uuid.UUID `gorm:"not null"`
	Tenant               Tenant    `gorm:"foreignKey:TenantID"`
	ReferenceID          string    `gorm:"type:varchar(255);not null"`
	PlannedDate          time.Time `gorm:"type:date;not null"`
}

func (p Plan) Map() domain.Plan {
	return domain.Plan{
		ReferenceID: p.ReferenceID,
		PlannedDate: p.PlannedDate,
	}
}
