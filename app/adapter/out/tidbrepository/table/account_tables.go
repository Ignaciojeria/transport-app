package table

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	ID int64 `gorm:"primaryKey"`
	// Contacto asociado a la orden
	ContactID int64   `gorm:"not null"`             // Clave foránea al Contact
	Contact   Contact `gorm:"foreignKey:ContactID"` // Relación con Contact
	IsActive  bool    `gorm:"not null;index"`
	// Dirección Origen Cuenta
	OriginNodeInfoID     int64                 `gorm:"default:null"` // ID del NodeInfo de origen
	OriginNodeInfo       NodeInfo              `gorm:"foreignKey:OriginNodeInfoID"`
	AccountOrganizations []AccountOrganization `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type AccountOrganization struct {
	gorm.Model
	ID             int64  `gorm:"primaryKey"`
	AccountID      int64  `gorm:"not null;uniqueIndex:idx_account_organization_country"`
	OrganizationID int64  `gorm:"not null;uniqueIndex:idx_account_organization_country"`
	Country        string `gorm:"type:varchar(5);not null;uniqueIndex:idx_account_organization_country"`
}
