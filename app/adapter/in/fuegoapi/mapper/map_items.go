package mapper

import "transport-app/app/domain"

func MapItemsToDomain(items []struct {
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
}) []domain.Item {
	mapped := make([]domain.Item, len(items))
	for i, item := range items {
		mapped[i] = domain.Item{
			Description: item.Description,
			Dimensions: domain.Dimensions{
				Height: item.Dimensions.Height,
				Width:  item.Dimensions.Width,
				Length: item.Dimensions.Length,
				Unit:   item.Dimensions.Unit,
			},
			Insurance: domain.Insurance{
				Currency:  item.Insurance.Currency,
				UnitValue: item.Insurance.UnitValue,
			},
			Skills: MapSkillsToDomain(item.Skills),
			Quantity: domain.Quantity{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
			Sku: item.Sku,
			Weight: domain.Weight{
				Unit:  item.Weight.Unit,
				Value: item.Weight.Value,
			},
		}
	}
	return mapped
}

func MapSkillsToDomain(skills []struct {
	Type        string `json:"type" example:"requires_refrigeration"`
	Value       string `json:"value" example:"2°C-8°C"`
	Description string `json:"description" example:"Producto refrigerado"`
}) []domain.Skill {
	mapped := make([]domain.Skill, len(skills))
	for i, skill := range skills {
		mapped[i] = domain.Skill{
			Type:        skill.Type,
			Value:       skill.Value,
			Description: skill.Description,
		}
	}
	return mapped
}

func MapSkillsFromDomain(skills []domain.Skill) []struct {
	Type        string `json:"type" example:"requires_refrigeration"`
	Value       string `json:"value" example:"2°C-8°C"`
	Description string `json:"description" example:"Producto refrigerado"`
} {
	mapped := make([]struct {
		Type        string `json:"type" example:"requires_refrigeration"`
		Value       string `json:"value" example:"2°C-8°C"`
		Description string `json:"description" example:"Producto refrigerado"`
	}, len(skills))
	for i, skill := range skills {
		mapped[i] = struct {
			Type        string `json:"type" example:"requires_refrigeration"`
			Value       string `json:"value" example:"2°C-8°C"`
			Description string `json:"description" example:"Producto refrigerado"`
		}{
			Type:        skill.Type,
			Value:       skill.Value,
			Description: skill.Description,
		}
	}
	return mapped
}

func MapItemsFromDomain(items []domain.Item) []struct {
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
} {
	mapped := make([]struct {
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
	}, len(items))

	for i, item := range items {
		mapped[i] = struct {
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
		}{
			Description: item.Description,
			Dimensions: struct {
				Length float64 `json:"length" example:"100.0"`
				Height float64 `json:"height" example:"100.0"`
				Unit   string  `json:"unit" example:"cm"`
				Width  float64 `json:"width" example:"100.0"`
			}{
				Length: item.Dimensions.Length,
				Height: item.Dimensions.Height,
				Unit:   item.Dimensions.Unit,
				Width:  item.Dimensions.Width,
			},
			Insurance: struct {
				Currency  string  `json:"currency" example:"CLP"`
				UnitValue float64 `json:"unitValue" example:"10000"`
			}{
				Currency:  item.Insurance.Currency,
				UnitValue: item.Insurance.UnitValue,
			},
			Skills: MapSkillsFromDomain(item.Skills),
			Quantity: struct {
				QuantityNumber int    `json:"quantityNumber" example:"1"`
				QuantityUnit   string `json:"quantityUnit" example:"unit"`
			}{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
			Sku: item.Sku,
			Weight: struct {
				Unit  string  `json:"unit" example:"kg"`
				Value float64 `json:"value" example:"1.8"`
			}{
				Unit:  item.Weight.Unit,
				Value: item.Weight.Value,
			},
		}
	}

	return mapped
}
