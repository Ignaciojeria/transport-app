package mapper

import "transport-app/app/domain"

// MapPackagesToDomain convierte estructuras anónimas de paquetes a domain.Package
func MapPackagesToDomain(packages []struct {
	SizeCategory string `json:"sizeCategory" example:"XL"`
	Dimensions   struct {
		Length float64 `json:"length" example:"100.0"`
		Height float64 `json:"height" example:"100.0"`
		Unit   string  `json:"unit" example:"cm"`
		Width  float64 `json:"width" example:"100.0"`
	} `json:"dimensions"`
	Insurance struct {
		Currency  string  `json:"currency" example:"CLP"`
		UnitValue float64 `json:"unitValue" example:"10000"`
	} `json:"insurance"`
	Items []struct {
		Description string `json:"description" example:"Cama 1 plaza"`
		Dimensions  struct {
			Length float64 `json:"length" example:"100.0"`
			Height float64 `json:"height" example:"100.0"`
			Unit   string  `json:"unit" example:"cm"`
			Width  float64 `json:"width" example:"100.0"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string  `json:"currency" example:"CLP"`
			UnitValue float64 `json:"unitValue" example:"10000"`
		} `json:"insurance"`
		Skills []struct {
			Type        string `json:"type" example:"requires_refrigeration"`
			Value       string `json:"value" example:"2°C-8°C"`
			Description string `json:"description" example:"Producto refrigerado"`
		} `json:"skills"`
		Quantity struct {
			QuantityNumber int    `json:"quantityNumber" example:"1"`
			QuantityUnit   string `json:"quantityUnit" example:"unit"`
		} `json:"quantity"`
		Sku    string `json:"sku" example:"1234567890"`
		Weight struct {
			Unit  string  `json:"unit" example:"kg"`
			Value float64 `json:"value" example:"1.8"`
		} `json:"weight"`
	} `json:"items"`
	Lpn    string `json:"lpn" example:"1234567890"`
	Labels []struct {
		Type  string `json:"type" example:"skill"`
		Value string `json:"value" example:"requires_refrigeration"`
	} `json:"labels"`
	Weight struct {
		Unit  string  `json:"unit" example:"kg"`
		Value float64 `json:"value" example:"1.8"`
	} `json:"weight"`
}) []domain.DeliveryUnit {
	mapped := make([]domain.DeliveryUnit, len(packages))
	for i, pkg := range packages {
		mapped[i] = domain.DeliveryUnit{
			SizeCategory: domain.SizeCategory{Code: pkg.SizeCategory},
			Lpn:          pkg.Lpn,
			Status:       domain.Status{Status: domain.StatusAvailable},
			Dimensions: domain.Dimensions{
				Height: pkg.Dimensions.Height,
				Width:  pkg.Dimensions.Width,
				Length: pkg.Dimensions.Length,
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
			Labels: MapLabelsToDomain(pkg.Labels),
			Items:  MapItemsToDomain(pkg.Items),
		}
	}
	return mapped
}

func MapPackagesFromDomain(packages []domain.DeliveryUnit) []struct {
	Dimensions struct {
		Length float64 `json:"length"`
		Height float64 `json:"height"`
		Unit   string  `json:"unit"`
		Width  float64 `json:"width"`
	} `json:"dimensions"`
	Insurance struct {
		Currency  string  `json:"currency"`
		UnitValue float64 `json:"unitValue"`
	} `json:"insurance"`
	Items []struct {
		Description string `json:"description"`
		Dimensions  struct {
			Length float64 `json:"length"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		} `json:"insurance"`
		Skills []struct {
			Type        string `json:"type"`
			Value       string `json:"value"`
			Description string `json:"description"`
		} `json:"skills"`
		Quantity struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		Sku    string `json:"sku"`
		Weight struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"weight"`
	} `json:"items"`
	Lpn    string `json:"lpn"`
	Weight struct {
		Unit  string  `json:"unit"`
		Value float64 `json:"value"`
	} `json:"weight"`
} {
	mapped := make([]struct {
		Dimensions struct {
			Length float64 `json:"length"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		} `json:"insurance"`
		Items []struct {
			Description string `json:"description"`
			Dimensions  struct {
				Length float64 `json:"length"`
				Height float64 `json:"height"`
				Unit   string  `json:"unit"`
				Width  float64 `json:"width"`
			} `json:"dimensions"`
			Insurance struct {
				Currency  string  `json:"currency"`
				UnitValue float64 `json:"unitValue"`
			} `json:"insurance"`
			Skills []struct {
				Type        string `json:"type"`
				Value       string `json:"value"`
				Description string `json:"description"`
			} `json:"skills"`
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			Sku    string `json:"sku"`
			Weight struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"weight"`
		} `json:"items"`
		Lpn    string `json:"lpn"`
		Weight struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"weight"`
	}, len(packages))

	for i, pkg := range packages {
		mapped[i] = struct {
			Dimensions struct {
				Length float64 `json:"length"`
				Height float64 `json:"height"`
				Unit   string  `json:"unit"`
				Width  float64 `json:"width"`
			} `json:"dimensions"`
			Insurance struct {
				Currency  string  `json:"currency"`
				UnitValue float64 `json:"unitValue"`
			} `json:"insurance"`
			Items []struct {
				Description string `json:"description"`
				Dimensions  struct {
					Length float64 `json:"length"`
					Height float64 `json:"height"`
					Unit   string  `json:"unit"`
					Width  float64 `json:"width"`
				} `json:"dimensions"`
				Insurance struct {
					Currency  string  `json:"currency"`
					UnitValue float64 `json:"unitValue"`
				} `json:"insurance"`
				Skills []struct {
					Type        string `json:"type"`
					Value       string `json:"value"`
					Description string `json:"description"`
				} `json:"skills"`
				Quantity struct {
					QuantityNumber int    `json:"quantityNumber"`
					QuantityUnit   string `json:"quantityUnit"`
				} `json:"quantity"`
				Sku    string `json:"sku"`
				Weight struct {
					Unit  string  `json:"unit"`
					Value float64 `json:"value"`
				} `json:"weight"`
			} `json:"items"`
			Lpn    string `json:"lpn"`
			Weight struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"weight"`
		}{
			Dimensions: struct {
				Length float64 `json:"length"`
				Height float64 `json:"height"`
				Unit   string  `json:"unit"`
				Width  float64 `json:"width"`
			}{
				Length: pkg.Dimensions.Length,
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
			Lpn: pkg.Lpn,
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

func MapLabelsToDomain(labels []struct {
	Type  string `json:"type" example:"skill"`
	Value string `json:"value" example:"requires_refrigeration"`
}) []domain.Reference {
	mapped := make([]domain.Reference, len(labels))
	for i, label := range labels {
		mapped[i] = domain.Reference{
			Type:  label.Type,
			Value: label.Value,
		}
	}
	return mapped
}

func MapLabelsFromDomain(labels []domain.Reference) []struct {
	Type  string `json:"type" example:"skill"`
	Value string `json:"value" example:"requires_refrigeration"`
} {
	mapped := make([]struct {
		Type  string `json:"type" example:"skill"`
		Value string `json:"value" example:"requires_refrigeration"`
	}, len(labels))
	for i, label := range labels {
		mapped[i] = struct {
			Type  string `json:"type" example:"skill"`
			Value string `json:"value" example:"requires_refrigeration"`
		}{
			Type:  label.Type,
			Value: label.Value,
		}
	}
	return mapped
}
