package domain

import "context"

type VehicleCategory struct {
	Type                string
	MaxPackagesQuantity int
}

func (vc VehicleCategory) DocID(ctx context.Context) DocumentID {
	return Hash(ctx, vc.Type)
}

func (vc VehicleCategory) UpdateIfChanged(in VehicleCategory) (VehicleCategory, bool) {
	changed := false

	if in.MaxPackagesQuantity != 0 && in.MaxPackagesQuantity != vc.MaxPackagesQuantity {
		vc.MaxPackagesQuantity = in.MaxPackagesQuantity
		changed = true
	}

	return vc, changed
}
