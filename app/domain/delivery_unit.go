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
	// Simplified fields (matching optimization domain)
	Volume          int64 // Volume in cm³ (by convention)
	WeightValue     int64 // Weight in grams (by convention)
	InsuranceValue  int64 // Insurance in CLP (by convention)
	// Legacy complex fields (for backward compatibility)
	Dimensions      Dimensions
	Weight          Weight
	Insurance       Insurance
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

// GetWeightValue returns the weight value, prioritizing WeightValue over Weight.Value
func (p DeliveryUnit) GetWeightValue() int64 {
	if p.WeightValue > 0 {
		return p.WeightValue
	}
	return p.Weight.Value
}

// GetInsuranceValue returns the insurance value, prioritizing InsuranceValue over Insurance.UnitValue
func (p DeliveryUnit) GetInsuranceValue() int64 {
	if p.InsuranceValue > 0 {
		return p.InsuranceValue
	}
	return p.Insurance.UnitValue
}

// GetVolume returns the volume value
func (p DeliveryUnit) GetVolume() int64 {
	return p.Volume
}

// SetSimpleValues sets the simplified values and syncs with legacy structures
func (p *DeliveryUnit) SetSimpleValues(volume, weight, insurance int64) {
	p.Volume = volume
	p.WeightValue = weight
	p.InsuranceValue = insurance
	
	// Sync with legacy structures for backward compatibility
	if weight > 0 {
		p.Weight = Weight{Value: weight, Unit: "g"}
	}
	if insurance > 0 {
		p.Insurance = Insurance{UnitValue: insurance, Currency: "CLP"}
	}
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
		Insurance: p.GetInsuranceValue(),
		Volume:    p.GetVolume(),
		Weight:    p.GetWeightValue(),
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
	
	deliveryUnit := DeliveryUnit{
		Lpn:   optDU.Lpn,
		Items: items,
	}
	
	// Set simplified values
	deliveryUnit.SetSimpleValues(optDU.Volume, optDU.Weight, optDU.Insurance)
	
	return deliveryUnit
}
