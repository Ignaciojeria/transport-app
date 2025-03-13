package mapper

import "transport-app/app/domain"

// MapPackagesToDomain convierte estructuras anónimas de paquetes a domain.Package
func MapPackagesToDomain(packages []struct {
	Dimensions struct {
		Depth  float64 `json:"depth"`
		Height float64 `json:"height"`
		Unit   string  `json:"unit"`
		Width  float64 `json:"width"`
	} `json:"dimensions"`
	Insurance struct {
		Currency  string  `json:"currency"`
		UnitValue float64 `json:"unitValue"`
	} `json:"insurance"`
	ItemReferences []struct {
		Quantity struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		ReferenceID string `json:"referenceId"`
	} `json:"itemReferences"`
	Lpn    string `json:"lpn"`
	Weight struct {
		Unit  string  `json:"unit"`
		Value float64 `json:"value"`
	} `json:"weight"`
}) []domain.Package {
	mapped := make([]domain.Package, len(packages))
	for i, pkg := range packages {
		mapped[i] = domain.Package{
			Lpn: pkg.Lpn,
			Dimensions: domain.Dimensions{
				Height: pkg.Dimensions.Height,
				Width:  pkg.Dimensions.Width,
				Depth:  pkg.Dimensions.Depth,
				Unit:   pkg.Dimensions.Unit,
			},
			Weight: domain.Weight{
				Unit:  pkg.Weight.Unit,
				Value: pkg.Weight.Value,
			},
			Insurance: domain.Insurance{
				Currency:  pkg.Insurance.Currency,
				UnitValue: pkg.Insurance.UnitValue,
			},
			ItemReferences: MapItemReferencesToDomain(pkg.ItemReferences),
		}
	}
	return mapped
}

func MapPackagesFromDomain(packages []domain.Package) []struct {
	Dimensions struct {
		Depth  float64 `json:"depth"`
		Height float64 `json:"height"`
		Unit   string  `json:"unit"`
		Width  float64 `json:"width"`
	} `json:"dimensions"`
	Insurance struct {
		Currency  string  `json:"currency"`
		UnitValue float64 `json:"unitValue"`
	} `json:"insurance"`
	ItemReferences []struct {
		Quantity struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		ReferenceID string `json:"referenceId"`
	} `json:"itemReferences"`
	Lpn    string `json:"lpn"`
	Weight struct {
		Unit  string  `json:"unit"`
		Value float64 `json:"value"`
	} `json:"weight"`
} {
	mapped := make([]struct {
		Dimensions struct {
			Depth  float64 `json:"depth"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		} `json:"insurance"`
		ItemReferences []struct {
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			ReferenceID string `json:"referenceId"`
		} `json:"itemReferences"`
		Lpn    string `json:"lpn"`
		Weight struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"weight"`
	}, len(packages))

	for i, pkg := range packages {
		mapped[i] = struct {
			Dimensions struct {
				Depth  float64 `json:"depth"`
				Height float64 `json:"height"`
				Unit   string  `json:"unit"`
				Width  float64 `json:"width"`
			} `json:"dimensions"`
			Insurance struct {
				Currency  string  `json:"currency"`
				UnitValue float64 `json:"unitValue"`
			} `json:"insurance"`
			ItemReferences []struct {
				Quantity struct {
					QuantityNumber int    `json:"quantityNumber"`
					QuantityUnit   string `json:"quantityUnit"`
				} `json:"quantity"`
				ReferenceID string `json:"referenceId"`
			} `json:"itemReferences"`
			Lpn    string `json:"lpn"`
			Weight struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"weight"`
		}{
			Dimensions: struct {
				Depth  float64 `json:"depth"`
				Height float64 `json:"height"`
				Unit   string  `json:"unit"`
				Width  float64 `json:"width"`
			}{
				Depth:  pkg.Dimensions.Depth,
				Height: pkg.Dimensions.Height,
				Unit:   pkg.Dimensions.Unit,
				Width:  pkg.Dimensions.Width,
			},
			Insurance: struct {
				Currency  string  `json:"currency"`
				UnitValue float64 `json:"unitValue"`
			}{
				Currency:  pkg.Insurance.Currency,
				UnitValue: pkg.Insurance.UnitValue,
			},
			ItemReferences: MapItemReferencesFromDomain(pkg.ItemReferences),
			Lpn:            pkg.Lpn,
			Weight: struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			}{
				Unit:  pkg.Weight.Unit,
				Value: pkg.Weight.Value,
			},
		}
	}

	return mapped
}
