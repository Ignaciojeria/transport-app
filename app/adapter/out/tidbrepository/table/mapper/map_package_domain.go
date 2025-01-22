package mapper

/*
import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapPackageDomain(pkg table.Package) domain.Package {
	return domain.Package{
		ID:             pkg.ID,
		Lpn:            pkg.Lpn,
		Dimensions:     mapTableDimensionsToDomain(pkg.JSONDimensions),
		Weight:         mapTableWeightToDomain(pkg.JSONWeight),
		Insurance:      mapTableInsuranceToDomain(pkg.JSONInsurance),
		ItemReferences: mapTableItemsToDomain(pkg.JSONItemsReferences),
	}
}

func mapTableDimensionsToDomain(dim table.JSONDimensions) domain.Dimensions {
	return domain.Dimensions{
		Height: dim.Height,
		Width:  dim.Width,
		Depth:  dim.Depth,
		Unit:   dim.Unit,
	}
}

func mapTableWeightToDomain(weight table.JSONWeight) domain.Weight {
	return domain.Weight{
		Value: weight.WeightValue,
		Unit:  weight.WeightUnit,
	}
}

func mapTableInsuranceToDomain(ins table.JSONInsurance) domain.Insurance {
	return domain.Insurance{
		UnitValue: ins.UnitValue,
		Currency:  ins.Currency,
	}
}

func mapTableItemsToDomain(items table.JSONItemReferences) []domain.ItemReference {
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
*/
