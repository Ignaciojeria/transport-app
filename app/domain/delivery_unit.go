package domain

import (
	"context"
	"fmt"
	"sort"
)

type DeliveryUnit struct {
	SizeCategory    SizeCategory
	Lpn             string
	Dimensions      Dimensions
	Weight          Weight
	Insurance       Insurance
	Index           int
	SkuIndex        string
	Status          Status
	ConfirmDelivery ConfirmDelivery
	Items           []Item
	Labels          []Reference `json:"labels"`
}

func (p DeliveryUnit) DocID(ctx context.Context, otherReference string) DocumentID {
	if p.Lpn != "" {
		return HashByTenant(ctx, p.Lpn)
	}

	var allInputs []string

	// Agregar referencia externa
	allInputs = append(allInputs, otherReference)

	// Agregar índice del paquete (por posición)
	allInputs = append(allInputs, fmt.Sprintf("index:%d", p.Index))

	// Agregar SKUs ordenados
	skus := make([]string, 0, len(p.Items))
	for _, item := range p.Items {
		skus = append(skus, item.Sku)
	}
	sort.Strings(skus)
	allInputs = append(allInputs, skus...)

	return HashByTenant(ctx, allInputs...)
}

func SearchPackageByLpn(pcks []DeliveryUnit, lpn string) DeliveryUnit {
	for _, pck := range pcks {
		if pck.Lpn == lpn {
			return pck
		}
	}
	return DeliveryUnit{}
}
func (p DeliveryUnit) UpdateIfChanged(newPackage DeliveryUnit) (DeliveryUnit, bool) {
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

/*
func (p *Package) ExplodeIfNoLpn() []Package {
	// Si el paquete tiene LPN, se considera agrupado y no se descompone
	if p.Lpn != "" {
		return []Package{*p}
	}

	var exploded []Package
	for _, item := range p.Items {
		exploded = append(exploded, Package{
			Dimensions: p.Dimensions, // Puedes replicar si aplica
			Weight:     item.Weight,  // O usar el general si prefieres
			Insurance:  item.Insurance,
			Items:      []Item{item},
		})
	}
	return exploded
}
*/
