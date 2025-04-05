package domain

type VehicleCategory struct {
	Organization        Organization
	Type                string
	MaxPackagesQuantity int
}

func (vc VehicleCategory) DocID() DocumentID {
	return Hash(vc.Organization, vc.Type)
}

func (vc VehicleCategory) UpdateIfChanged(in VehicleCategory) VehicleCategory {
	if in.MaxPackagesQuantity != 0 {
		vc.MaxPackagesQuantity = in.MaxPackagesQuantity
	}
	return vc
}
