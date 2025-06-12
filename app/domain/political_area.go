package domain

import "context"

type PoliticalArea struct {
	Code       string
	Province   string
	State      string
	District   string
	TimeZone   string
	Confidence CoordinatesConfidence
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
