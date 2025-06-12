package domain

import (
	"context"
	"transport-app/app/shared/utils"

	"github.com/paulmach/orb"
)

type CoordinatesConfidence struct {
	Level   float64
	Message string
	Reason  string
}

type Coordinates struct {
	Point      orb.Point
	Source     string
	Confidence CoordinatesConfidence
}

type AddressInfo struct {
	Contact       Contact
	PoliticalArea PoliticalArea
	AddressLine1  string
	AddressLine2  string
	Coordinates   Coordinates
	ZipCode       string
}

func (a AddressInfo) DocID(ctx context.Context) DocumentID {
	return HashByTenant(
		ctx,
		a.AddressLine1,
		a.AddressLine2,
		a.PoliticalArea.District,
		a.PoliticalArea.Province,
		a.PoliticalArea.State)
}

func (a *AddressInfo) UpdatePoint(point orb.Point) {
	a.Coordinates.Point = point
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
	a.PoliticalArea.State = utils.NormalizeText(a.PoliticalArea.State)
	a.PoliticalArea.Province = utils.NormalizeText(a.PoliticalArea.Province)
	a.PoliticalArea.District = utils.NormalizeText(a.PoliticalArea.District)
	a.PoliticalArea.TimeZone = utils.NormalizeText(a.PoliticalArea.TimeZone)
}

func (a AddressInfo) UpdateIfChanged(newAddress AddressInfo) (AddressInfo, bool) {
	updated := a
	changed := false

	if newAddress.AddressLine1 != "" && newAddress.AddressLine1 != a.AddressLine1 {
		updated.AddressLine1 = newAddress.AddressLine1
		changed = true
	}
	if newAddress.Coordinates.Point[1] != 0 && newAddress.Coordinates.Point[1] != a.Coordinates.Point[1] {
		updated.Coordinates.Point[1] = newAddress.Coordinates.Point[1]
		changed = true
	}
	if newAddress.Coordinates.Point[0] != 0 && newAddress.Coordinates.Point[0] != a.Coordinates.Point[0] {
		updated.Coordinates.Point[0] = newAddress.Coordinates.Point[0]
		changed = true
	}

	// Validar cambios en el nivel de confianza
	if newAddress.Coordinates.Confidence.Level != a.Coordinates.Confidence.Level {
		updated.Coordinates.Confidence.Level = newAddress.Coordinates.Confidence.Level
		changed = true
	}
	if newAddress.Coordinates.Confidence.Message != "" && newAddress.Coordinates.Confidence.Message != a.Coordinates.Confidence.Message {
		updated.Coordinates.Confidence.Message = newAddress.Coordinates.Confidence.Message
		changed = true
	}
	if newAddress.Coordinates.Confidence.Reason != "" && newAddress.Coordinates.Confidence.Reason != a.Coordinates.Confidence.Reason {
		updated.Coordinates.Confidence.Reason = newAddress.Coordinates.Confidence.Reason
		changed = true
	}

	if newAddress.PoliticalArea.State != "" && newAddress.PoliticalArea.State != a.PoliticalArea.State {
		updated.PoliticalArea.State = newAddress.PoliticalArea.State
		changed = true
	}

	if newAddress.PoliticalArea.Province != "" && newAddress.PoliticalArea.Province != a.PoliticalArea.Province {
		updated.PoliticalArea.Province = newAddress.PoliticalArea.Province
		changed = true
	}

	if newAddress.PoliticalArea.District != "" && newAddress.PoliticalArea.District != a.PoliticalArea.District {
		updated.PoliticalArea.District = newAddress.PoliticalArea.District
		changed = true
	}

	if newAddress.ZipCode != "" && newAddress.ZipCode != a.ZipCode {
		updated.ZipCode = newAddress.ZipCode
		changed = true
	}
	if newAddress.PoliticalArea.TimeZone != "" && newAddress.PoliticalArea.TimeZone != a.PoliticalArea.TimeZone {
		updated.PoliticalArea.TimeZone = newAddress.PoliticalArea.TimeZone
		changed = true
	}

	// Validar cambios en el nivel de confianza del PoliticalArea
	if newAddress.PoliticalArea.Confidence.Level != a.PoliticalArea.Confidence.Level {
		updated.PoliticalArea.Confidence.Level = newAddress.PoliticalArea.Confidence.Level
		changed = true
	}
	if newAddress.PoliticalArea.Confidence.Message != "" && newAddress.PoliticalArea.Confidence.Message != a.PoliticalArea.Confidence.Message {
		updated.PoliticalArea.Confidence.Message = newAddress.PoliticalArea.Confidence.Message
		changed = true
	}
	if newAddress.PoliticalArea.Confidence.Reason != "" && newAddress.PoliticalArea.Confidence.Reason != a.PoliticalArea.Confidence.Reason {
		updated.PoliticalArea.Confidence.Reason = newAddress.PoliticalArea.Confidence.Reason
		changed = true
	}

	return updated, changed
}

func (a AddressInfo) FullAddress() string {
	parts := []string{
		a.AddressLine1,
		a.PoliticalArea.District,
		a.PoliticalArea.Province,
		a.PoliticalArea.State,
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

func (a AddressInfo) Equals(ctx context.Context, other AddressInfo) bool {
	return a.DocID(ctx) == other.DocID(ctx)
}
