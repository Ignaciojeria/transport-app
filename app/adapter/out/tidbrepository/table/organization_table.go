package table

import "gorm.io/gorm"

type ApiKey struct {
	gorm.Model
	ID             int64        `gorm:"primaryKey"`
	OrganizationID int64        `gorm:"not null;index"`            // ID de la organización asociada
	Organization   Organization `gorm:"foreignKey:OrganizationID"` // Relación con la tabla Organization
	Key            string       `gorm:"not null;unique"`           // Clave única
	Status         string       `gorm:"default:active"`            // Estado: activo, revocado, etc.
}

type Organization struct {
	gorm.Model
	ID        int64                 `gorm:"primaryKey"`
	ApiKeyID  int64                 `gorm:"not null;index"` // Relación con ApiKey
	Email     string                `gorm:"type:varchar(255);not null;unique"`
	Name      string                `gorm:"type:varchar(255);not null;"`
	Countries []OrganizationCountry `gorm:"foreignKey:OrganizationID"` // Relación con países
}

type OrganizationCountry struct {
	gorm.Model
	ID             int64  `gorm:"primaryKey"`
	OrganizationID int64  `gorm:"not null;uniqueIndex:idx_organization_country"`              // Parte del índice compuesto
	Country        string `gorm:"type:char(5);not null;uniqueIndex:idx_organization_country"` // Código ISO de 2 caracteres
}
