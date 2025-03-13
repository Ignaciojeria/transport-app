package mapper

import "transport-app/app/domain"

func MapItemsToDomain(items []struct {
	Description string `json:"description"`
	Dimensions  struct {
		Depth  float64 `json:"depth"`
		Height float64 `json:"height"`
		Unit   string  `json:"unit"`
		Width  float64 `json:"width"`
	} `json:"dimensions"`
	Insurance struct {
		Currency  string  `json:"currency"`
		UnitValue float64 `json:"unitValue"`
	} `json:"insurance"`
	LogisticCondition string `json:"logisticCondition"`
	Quantity          struct {
		QuantityNumber int    `json:"quantityNumber"`
		QuantityUnit   string `json:"quantityUnit"`
	} `json:"quantity"`
	ReferenceID string `json:"referenceId"`
	Weight      struct {
		Unit  string  `json:"unit"`
		Value float64 `json:"value"`
	} `json:"weight"`
}) []domain.Item {
	mapped := make([]domain.Item, len(items))
	for i, item := range items {
		mapped[i] = domain.Item{
			ReferenceID:       domain.ReferenceID(item.ReferenceID),
			LogisticCondition: item.LogisticCondition,
			Quantity: domain.Quantity{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
			Insurance: domain.Insurance{
				Currency:  item.Insurance.Currency,
				UnitValue: item.Insurance.UnitValue,
			},
			Description: item.Description,
			Dimensions: domain.Dimensions{
				Height: item.Dimensions.Height,
				Width:  item.Dimensions.Width,
				Depth:  item.Dimensions.Depth,
				Unit:   item.Dimensions.Unit,
			},
			Weight: domain.Weight{
				Unit:  item.Weight.Unit,
				Value: item.Weight.Value,
			},
		}
	}
	return mapped
}

func MapItemsFromDomain(items []domain.Item) []struct {
	Description string `json:"description"`
	Dimensions  struct {
		Depth  float64 `json:"depth"`
		Height float64 `json:"height"`
		Unit   string  `json:"unit"`
		Width  float64 `json:"width"`
	} `json:"dimensions"`
	Insurance struct {
		Currency  string  `json:"currency"`
		UnitValue float64 `json:"unitValue"`
	} `json:"insurance"`
	LogisticCondition string `json:"logisticCondition"`
	Quantity          struct {
		QuantityNumber int    `json:"quantityNumber"`
		QuantityUnit   string `json:"quantityUnit"`
	} `json:"quantity"`
	ReferenceID string `json:"referenceId"`
	Weight      struct {
		Unit  string  `json:"unit"`
		Value float64 `json:"value"`
	} `json:"weight"`
} {
	mapped := make([]struct {
		Description string `json:"description"`
		Dimensions  struct {
			Depth  float64 `json:"depth"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		} `json:"insurance"`
		LogisticCondition string `json:"logisticCondition"`
		Quantity          struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		ReferenceID string `json:"referenceId"`
		Weight      struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"weight"`
	}, len(items))

	for i, item := range items {
		mapped[i] = struct {
			Description string `json:"description"`
			Dimensions  struct {
				Depth  float64 `json:"depth"`
				Height float64 `json:"height"`
				Unit   string  `json:"unit"`
				Width  float64 `json:"width"`
			} `json:"dimensions"`
			Insurance struct {
				Currency  string  `json:"currency"`
				UnitValue float64 `json:"unitValue"`
			} `json:"insurance"`
			LogisticCondition string `json:"logisticCondition"`
			Quantity          struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			ReferenceID string `json:"referenceId"`
			Weight      struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"weight"`
		}{
			Description: item.Description,
			Dimensions: struct {
				Depth  float64 `json:"depth"`
				Height float64 `json:"height"`
				Unit   string  `json:"unit"`
				Width  float64 `json:"width"`
			}{
				Depth:  item.Dimensions.Depth,
				Height: item.Dimensions.Height,
				Unit:   item.Dimensions.Unit,
				Width:  item.Dimensions.Width,
			},
			Insurance: struct {
				Currency  string  `json:"currency"`
				UnitValue float64 `json:"unitValue"`
			}{
				Currency:  item.Insurance.Currency,
				UnitValue: item.Insurance.UnitValue,
			},
			LogisticCondition: item.LogisticCondition,
			Quantity: struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			}{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
			ReferenceID: string(item.ReferenceID),
			Weight: struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			}{
				Unit:  item.Weight.Unit,
				Value: item.Weight.Value,
			},
		}
	}

	return mapped
}
