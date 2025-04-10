package domain

import "context"

type VehicleCategory struct {
	Type                string
	MaxPackagesQuantity int
}

func (vc VehicleCategory) DocID(ctx context.Context) DocumentID {
	return Hash(ctx, vc.Type)
}

func (vc VehicleCategory) UpdateIfChanged(in VehicleCategory) VehicleCategory {
	if in.MaxPackagesQuantity != 0 {
		vc.MaxPackagesQuantity = in.MaxPackagesQuantity
	}
	return vc
}
