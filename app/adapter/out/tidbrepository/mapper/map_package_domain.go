package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapPackageDomain(pkg table.Package) domain.Package {
	return domain.Package{
		Lpn:            pkg.Lpn,
		Dimensions:     mapTableDimensionsToDomain(pkg.Dimensions),
		Weight:         mapTableWeightToDomain(pkg.Weight),
		Insurance:      mapTableInsuranceToDomain(pkg.Insurance),
		ItemReferences: mapTableItemsToDomain(pkg.JSONItems),
	}
}

func mapTableDimensionsToDomain(dim table.Dimensions) domain.Dimensions {
	return domain.Dimensions{
		Height: dim.Height,
		Width:  dim.Width,
		Depth:  dim.Depth,
		Unit:   dim.Unit,
	}
}

func mapTableWeightToDomain(weight table.Weight) domain.Weight {
	return domain.Weight{
		Value: weight.Value,
		Unit:  weight.Unit,
	}
}

func mapTableInsuranceToDomain(ins table.Insurance) domain.Insurance {
	return domain.Insurance{
		UnitValue: ins.UnitValue,
		Currency:  ins.Currency,
	}
}

func mapTableItemsToDomain(items table.JSONItems) []domain.ItemReference {
	mapped := make([]domain.ItemReference, len(items))
	for i, item := range items {
		mapped[i] = domain.ItemReference{
			ReferenceID: domain.ReferenceID(item.ReferenceID),
			Quantity: domain.Quantity{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
		}
	}
	return mapped
}
