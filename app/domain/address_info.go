package domain

import (
	"transport-app/app/shared/utils"

	"github.com/paulmach/orb"
)

type AddressInfo struct {
	ID                int64
	Organization      Organization
	Contact           Contact   `json:"contact"`
	State             string    `json:"state"`
	Province          string    `json:"province"`
	District          string    `json:"district"`
	AddressLine1      string    `json:"addressLine1"`
	Location          orb.Point // Punto original
	CorrectedLocation orb.Point // Punto corregido del plan (snapped)
	CorrectedDistance float64   // Distancia aplicada al punto corregido del plan
	ZipCode           string    `json:"zipCode"`
	TimeZone          string    `json:"timeZone"`
}

func (a AddressInfo) DocID() string {
	return Hash(
		a.Organization,
		a.AddressLine1,
		a.District,
		a.Province,
		a.State)
}

// Normalize limpia y formatea los valores de State, Province y District.
func (a *AddressInfo) Normalize() {
	a.AddressLine1 = utils.NormalizeText(a.AddressLine1)
	a.State = utils.NormalizeText(a.State)
	a.Province = utils.NormalizeText(a.Province)
	a.District = utils.NormalizeText(a.District)
}

func (a AddressInfo) UpdateIfChanged(newAddress AddressInfo) (AddressInfo, bool) {
	updated := a
	changed := false

	if newAddress.AddressLine1 != "" && newAddress.AddressLine1 != a.AddressLine1 {
		updated.AddressLine1 = newAddress.AddressLine1
		changed = true
	}
	if newAddress.Location[1] != 0 && newAddress.Location[1] != a.Location[1] {
		updated.Location[1] = newAddress.Location[1]
		changed = true
	}
	if newAddress.Location[0] != 0 && newAddress.Location[0] != a.Location[0] {
		updated.Location[0] = newAddress.Location[0]
		changed = true
	}
	if newAddress.State != "" && newAddress.State != a.State {
		updated.State = newAddress.State
		changed = true
	}

	if newAddress.Province != "" && newAddress.Province != a.Province {
		updated.Province = newAddress.Province
		changed = true
	}
	if newAddress.District != "" && newAddress.District != a.District {
		updated.District = newAddress.District
		changed = true
	}
	if newAddress.ZipCode != "" && newAddress.ZipCode != a.ZipCode {
		updated.ZipCode = newAddress.ZipCode
		changed = true
	}
	if newAddress.TimeZone != "" && newAddress.TimeZone != a.TimeZone {
		updated.TimeZone = newAddress.TimeZone
		changed = true
	}

	return updated, changed
}

func (a AddressInfo) FullAddress() string {
	parts := []string{
		a.AddressLine1,
		a.District,
		a.Province,
		a.State,
		a.ZipCode,
	}
	return concatenateWithCommas(parts...)
}
