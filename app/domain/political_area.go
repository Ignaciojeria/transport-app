package domain

import "context"

type PoliticalArea struct {
	Code     string `gorm:"type:varchar(191);not null"`
	Province string `gorm:"type:varchar(191);not null"`
	State    string `gorm:"type:varchar(191);not null"`
	District string `gorm:"type:varchar(191);not null"`
	TimeZone string `gorm:"type:varchar(191);not null"`
}

func (p PoliticalArea) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, p.Province, p.State, p.District, p.Code)
}

func (p PoliticalArea) UpdateIfChanged(new PoliticalArea) (PoliticalArea, bool) {
	if p.TimeZone != new.TimeZone {
		return new, true
	}
	return new, false
}
