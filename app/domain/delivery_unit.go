package domain

import (
	"context"
	"sort"
	"transport-app/app/domain/optimization"
)

type DeliveryUnit struct {
	SizeCategory    SizeCategory
	Lpn             string
	noLPNReference  string
	Volume          *int64
	Weight          *int64
	Insurance       *int64
	Status          Status
	ConfirmDelivery ConfirmDelivery
	Items           []Item
	Labels          []Reference `json:"labels"`
	Skills          []Skill
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

	// Actualizar volumen si el puntero no es nil y el valor es diferente
	if newPackage.Volume != nil && (p.Volume == nil || *newPackage.Volume != *p.Volume) {
		p.Volume = newPackage.Volume
		changed = true
	}
	// Actualizar peso si el puntero no es nil y el valor es diferente
	if newPackage.Weight != nil && (p.Weight == nil || *newPackage.Weight != *p.Weight) {
		p.Weight = newPackage.Weight
		changed = true
	}
	// Actualizar seguro si el puntero no es nil y el valor es diferente
	if newPackage.Insurance != nil && (p.Insurance == nil || *newPackage.Insurance != *p.Insurance) {
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

// SetValues sets the simplified values directly
func (p *DeliveryUnit) SetValues(volume, weight, insurance int64) {
	p.Volume = &volume
	p.Weight = &weight
	p.Insurance = &insurance
}

// ToOptimizationDeliveryUnit converts this DeliveryUnit to the optimization domain structure
func (p DeliveryUnit) ToOptimizationDeliveryUnit() optimization.DeliveryUnit {
	items := make([]optimization.Item, len(p.Items))
	for i, item := range p.Items {
		items[i] = optimization.Item{
			Sku: item.Sku,
		}
	}

	return optimization.DeliveryUnit{
		Items:     items,
		Insurance: *p.Insurance,
		Volume:    *p.Volume,
		Weight:    *p.Weight,
		Lpn:       p.Lpn,
	}
}

// FromOptimizationDeliveryUnit creates a DeliveryUnit from the optimization domain structure
func FromOptimizationDeliveryUnit(optDU optimization.DeliveryUnit) DeliveryUnit {
	items := make([]Item, len(optDU.Items))
	for i, item := range optDU.Items {
		items[i] = Item{
			Sku: item.Sku,
		}
	}

	return DeliveryUnit{
		Lpn:       optDU.Lpn,
		Volume:    &optDU.Volume,
		Weight:    &optDU.Weight,
		Insurance: &optDU.Insurance,
		Items:     items,
	}
}
