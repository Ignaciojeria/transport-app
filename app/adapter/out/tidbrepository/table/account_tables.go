package table

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	ID                   int64                 `gorm:"primaryKey"`
	Email                string                `gorm:"type:varchar(191);not null;uniqueIndex"` // Cambia a varchar(191)
	NationalID           string                `gorm:"type:varchar(191);not null;index"`       // Asegúrate de que NationalID también tenga un tipo indexable
	OriginID             *int64                `gorm:"default:null"`                           // Cambia a puntero para permitir valores nulos
	IsActive             bool                  `gorm:"not null;index"`
	Origin               Origin                `gorm:"foreignKey:OriginID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	AccountOrganizations []AccountOrganization `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type AccountOrganization struct {
	gorm.Model
	ID             int64  `gorm:"primaryKey"`
	AccountID      int64  `gorm:"not null;uniqueIndex:idx_account_organization_country"`
	OrganizationID int64  `gorm:"not null;uniqueIndex:idx_account_organization_country"`
	Country        string `gorm:"type:varchar(5);not null;uniqueIndex:idx_account_organization_country"`
}
