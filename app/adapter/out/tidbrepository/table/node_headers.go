package table

import (
	"transport-app/app/domain"

	"github.com/google/uuid"
)

type NodeHeaders struct {
	ID       int64     `gorm:"primaryKey"`
	Commerce string    `gorm:"uniqueIndex:idx_commerce_consumer_org_country,length:50"`
	Consumer string    `gorm:"uniqueIndex:idx_commerce_consumer_org_country,length:50"`
	TenantID uuid.UUID `gorm:"not null;index;uniqueIndex:idx_commerce_consumer_org_country"`
	Tenant   Tenant    `gorm:"foreignKey:TenantID"`
}

func (m NodeHeaders) Map() domain.Headers {
	return domain.Headers{
		Consumer: m.Consumer,
		Commerce: m.Commerce,
	}
}
