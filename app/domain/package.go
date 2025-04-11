package domain

import "context"

type Package struct {
	Lpn        string
	Dimensions Dimensions
	Weight     Weight
	Insurance  Insurance
	Items      []Item
}

func (p Package) DocID(ctx context.Context, otherReference string) DocumentID {
	// Si el LPN existe, usarlo para generar el hash
	if p.Lpn != "" {
		return Hash(ctx, p.Lpn)
	}

	// Si no hay LPN, concatenar las referencias externas con los SKUs de los items
	var allInputs []string

	// Primero agregamos las referencias externas
	allInputs = append(allInputs, otherReference)

	// Luego agregamos los SKUs de los items
	for _, item := range p.Items {
		allInputs = append(allInputs, item.Sku)
	}

	// Generamos el hash con todos los inputs concatenados
	return Hash(ctx, allInputs...)
}

func SearchPackageByLpn(pcks []Package, lpn string) Package {
	for _, pck := range pcks {
		if pck.Lpn == lpn {
			return pck
		}
	}
	return Package{}
}
func (p Package) UpdateIfChanged(newPackage Package) (Package, bool) {
	changed := false

	// Actualizar Lpn
	if newPackage.Lpn != "" && newPackage.Lpn != p.Lpn {
		p.Lpn = newPackage.Lpn
		changed = true
	}

	// Actualizar dimensiones si no están vacías
	if newPackage.Dimensions != (Dimensions{}) {
		p.Dimensions = newPackage.Dimensions
		changed = true
	}

	// Actualizar peso si no está vacío
	if newPackage.Weight != (Weight{}) {
		p.Weight = newPackage.Weight
		changed = true
	}

	// Actualizar seguro si no está vacío
	if newPackage.Insurance != (Insurance{}) {
		p.Insurance = newPackage.Insurance
		changed = true
	}

	// Actualizar referencias de ítems
	if len(newPackage.Items) > 0 {
		p.Items = newPackage.Items
		changed = true
	}

	return p, changed
}
