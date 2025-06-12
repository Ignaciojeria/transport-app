package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PoliticalArea struct {
	ID         uint           `gorm:"primarykey"`
	Code       string         `gorm:"type:varchar(191);not null"`
	Province   string         `gorm:"type:varchar(191);not null"`
	State      string         `gorm:"type:varchar(191);not null"`
	District   string         `gorm:"type:varchar(191);not null"`
	ZipCode    string         `gorm:"type:varchar(191);not null"`
	TimeZone   string         `gorm:"type:varchar(191);not null"`
	DocumentID string         `gorm:"type:varchar(191);not null;uniqueIndex:idx_political_areas_document_id_tenant_id"`
	TenantID   uuid.UUID      `gorm:"type:char(36);not null;uniqueIndex:idx_political_areas_document_id_tenant_id"`
	CreatedAt  int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt  int64          `gorm:"autoUpdateTime:milli"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (pa PoliticalArea) Map() domain.PoliticalArea {
	return domain.PoliticalArea{
		Code:     pa.Code,
		Province: pa.Province,
		State:    pa.State,
		District: pa.District,
		TimeZone: pa.TimeZone,
	}
}

func (pa PoliticalArea) UpdateIfChanged(new domain.PoliticalArea) (domain.PoliticalArea, bool) {
	if pa.TimeZone != new.TimeZone {
		return new, true
	}
	return new, false
}
