package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PoliticalArea struct {
	ID              uint           `gorm:"primarykey"`
	Code            string         `gorm:"type:varchar(191);not null"`
	AdminAreaLevel1 string         `gorm:"type:varchar(191);not null"`
	AdminAreaLevel2 string         `gorm:"type:varchar(191);not null"`
	AdminAreaLevel3 string         `gorm:"type:varchar(191);not null"`
	AdminAreaLevel4 string         `gorm:"type:varchar(191);not null"`
	TimeZone        string         `gorm:"type:varchar(191);not null"`
	DocumentID      string         `gorm:"type:varchar(191);not null;uniqueIndex:idx_political_areas_document_id_tenant_id"`
	TenantID        uuid.UUID      `gorm:"type:char(36);not null;uniqueIndex:idx_political_areas_document_id_tenant_id"`
	CreatedAt       int64          `gorm:"autoCreateTime:milli"`
	UpdatedAt       int64          `gorm:"autoUpdateTime:milli"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (pa PoliticalArea) Map() domain.PoliticalArea {
	return domain.PoliticalArea{
		Code:            pa.Code,
		AdminAreaLevel1: pa.AdminAreaLevel1,
		AdminAreaLevel2: pa.AdminAreaLevel2,
		AdminAreaLevel3: pa.AdminAreaLevel3,
		AdminAreaLevel4: pa.AdminAreaLevel4,
		TimeZone:        pa.TimeZone,
	}
}

func (pa PoliticalArea) UpdateIfChanged(new domain.PoliticalArea) (domain.PoliticalArea, bool) {
	if pa.TimeZone != new.TimeZone {
		return new, true
	}
	return new, false
}
