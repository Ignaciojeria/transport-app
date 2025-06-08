package domain

import (
	"context"
	"sort"
)

type DeliveryUnit struct {
	SizeCategory    SizeCategory
	Lpn             string
	noLPNReference  string
	Dimensions      Dimensions
	Weight          Weight
	Insurance       Insurance
	Status          Status
	ConfirmDelivery ConfirmDelivery
	Items           []Item
	Labels          []Reference `json:"labels"`
}

func (p DeliveryUnit) DocID(ctx context.Context) DocumentID {
	if p.Lpn != "" {
		return HashByTenant(ctx, p.Lpn)
	}

	var allInputs []string

	// Primero agregar la referencia
	allInputs = append(allInputs, p.noLPNReference)

	// Luego agregar SKUs ordenados
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

func (p *DeliveryUnit) UpdateStatusBasedOnNonDelivery() {
	if p.ConfirmDelivery.NonDeliveryReason.IsEmpty() {
		p.Status = Status{
			Status: StatusFinished,
		}
	} else {
		p.Status = Status{
			Status: StatusPending,
		}
	}
}
