package domain

import (
	"strings"
	"transport-app/app/shared/utils"

	"github.com/paulmach/orb"
)

type AddressInfo struct {
	ID                int64
	Organization      Organization
	Contact           Contact   `json:"contact"`
	State             string    `json:"state"`
	Locality          string    `json:"locality"`
	Province          string    `json:"province"`
	District          string    `json:"district"`
	ProviderAddress   string    `json:"providerAddress"`
	AddressLine1      string    `json:"addressLine1"`
	AddressLine2      string    `json:"addressLine2"`
	AddressLine3      string    `json:"addressLine3"`
	Location          orb.Point // Punto original
	CorrectedLocation orb.Point // Punto corregido del plan (snapped)
	CorrectedDistance float64   // Distancia aplicada al punto corregido del plan
	ZipCode           string    `json:"zipCode"`
	TimeZone          string    `json:"timeZone"`
}

// Normalize limpia y formatea los valores de State, Province y District.
func (a *AddressInfo) Normalize() {
	a.State = utils.NormalizeText(a.State)
	a.ProviderAddress = strings.ToLower(a.ProviderAddress)
	a.AddressLine1 = strings.ToLower(a.AddressLine1)
	a.AddressLine2 = strings.ToLower(a.AddressLine2)
	a.Province = utils.NormalizeText(a.Province)
	a.District = utils.NormalizeText(a.District)
}

func (a AddressInfo) UpdateIfChanged(newAddress AddressInfo) AddressInfo {
	if newAddress.ID != 0 {
		a.ID = newAddress.ID
	}
	if newAddress.AddressLine1 != "" {
		a.AddressLine1 = newAddress.AddressLine1
	}
	if newAddress.AddressLine2 != "" {
		a.AddressLine2 = newAddress.AddressLine2
	}
	if newAddress.AddressLine3 != "" {
		a.AddressLine3 = newAddress.AddressLine3
	}
	if newAddress.Location[1] != 0 { // Si la latitud no es 0
		a.Location[1] = newAddress.Location[1]
	}
	if newAddress.Location[0] != 0 { // Si la longitud no es 0
		a.Location[0] = newAddress.Location[0]
	}
	if newAddress.State != "" {
		a.State = newAddress.State
	}
	if newAddress.Locality != "" {
		a.Locality = newAddress.Locality
	}
	if newAddress.Province != "" {
		a.Province = newAddress.Province
	}
	if newAddress.District != "" {
		a.District = newAddress.District
	}
	if newAddress.ZipCode != "" {
		a.ZipCode = newAddress.ZipCode
	}
	if newAddress.TimeZone != "" {
		a.TimeZone = newAddress.TimeZone
	}
	return a
}

func (addr AddressInfo) RawAddress() string {
	return concatenateWithCommas(addr.AddressLine1, addr.AddressLine2, addr.AddressLine3)
}
