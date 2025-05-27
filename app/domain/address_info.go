package domain

import (
	"context"
	"transport-app/app/shared/utils"

	"github.com/paulmach/orb"
)

type CoordinateSource string

type AddressInfo struct {
	Contact              Contact
	State                State
	Province             Province
	District             District
	AddressLine1         string
	AddressLine2         string
	Location             orb.Point
	RequiresManualReview bool
	CoordinateSource     string
	ZipCode              string
	TimeZone             string
}

func (a AddressInfo) DocID(ctx context.Context) DocumentID {
	return HashByTenant(
		ctx,
		a.AddressLine1,
		a.AddressLine2,
		a.District.String(),
		a.Province.String(),
		a.State.String())
}

func (a *AddressInfo) UpdatePoint(point orb.Point) {
	a.Location = point
}

func (a *AddressInfo) NormalizeAndGeocode(
	ctx context.Context,
	geocodeFn func(context.Context, AddressInfo) (orb.Point, error),
) error {
	a.ToLowerAndRemovePunctuation()
	point, err := geocodeFn(ctx, *a)
	if err != nil {
		return err
	}
	a.UpdatePoint(point)
	return nil
}

// Normalize limpia y formatea los valores de State, Province y District.
func (a *AddressInfo) ToLowerAndRemovePunctuation() {
	a.AddressLine1 = utils.NormalizeText(a.AddressLine1)
	a.State = State(utils.NormalizeText(a.State.String()))
	a.Province = Province(utils.NormalizeText(a.Province.String()))
	a.District = District(utils.NormalizeText(a.District.String()))
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

	if !newAddress.State.IsEmpty() && !newAddress.State.Equals(a.State) {
		updated.State = newAddress.State
		changed = true
	}

	if !newAddress.Province.IsEmpty() && !newAddress.Province.Equals(a.Province) {
		updated.Province = newAddress.Province
		changed = true
	}

	if !newAddress.District.IsEmpty() && !newAddress.District.Equals(a.District) {
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
		a.District.String(),
		a.Province.String(),
		a.State.String(),
		a.ZipCode,
	}
	return concatenateWithCommas(parts...)
}

func concatenateWithCommas(values ...string) string {
	result := ""
	for _, value := range values {
		if value != "" {
			if result != "" {
				result += ", "
			}
			result += value
		}
	}
	return result
}
