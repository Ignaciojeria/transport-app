package domain

type Package struct {
	Lpn            string
	Organization   Organization
	Dimensions     Dimensions
	Weight         Weight
	Insurance      Insurance
	ItemReferences []ItemReference
}

func (p Package) DocID() DocumentID {
	return Hash(p.Organization, p.Lpn)
}

func SearchPackageByLpn(pcks []Package, lpn string) Package {
	for _, pck := range pcks {
		if pck.Lpn == lpn {
			return pck
		}
	}
	return Package{}
}

func (p Package) UpdateIfChanged(newPackage Package) Package {
	// Actualizar Lpn
	if newPackage.Lpn != "" {
		p.Lpn = newPackage.Lpn
	}

	// Actualizar dimensiones si no están vacías
	if newPackage.Dimensions != (Dimensions{}) {
		p.Dimensions = newPackage.Dimensions
	}

	// Actualizar peso si no está vacío
	if newPackage.Weight != (Weight{}) {
		p.Weight = newPackage.Weight
	}

	// Actualizar seguro si no está vacío
	if newPackage.Insurance != (Insurance{}) {
		p.Insurance = newPackage.Insurance
	}

	// Actualizar referencias de ítems
	if len(newPackage.ItemReferences) > 0 {
		p.ItemReferences = newPackage.ItemReferences
	}
	return p
}
