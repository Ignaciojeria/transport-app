package domain

import "context"

type PoliticalArea struct {
	Code            string
	AdminAreaLevel1 string
	AdminAreaLevel2 string
	AdminAreaLevel3 string
	AdminAreaLevel4 string
	TimeZone        string
	Confidence      CoordinatesConfidence
}

func (p PoliticalArea) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, p.AdminAreaLevel1, p.AdminAreaLevel2, p.AdminAreaLevel3, p.Code)
}

func (p PoliticalArea) UpdateIfChanged(new PoliticalArea) (PoliticalArea, bool) {
	if p.TimeZone != new.TimeZone {
		return new, true
	}
	return new, false
}
