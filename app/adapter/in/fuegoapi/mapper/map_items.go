package mapper

import "transport-app/app/domain"

func MapItemsToDomain(items []struct {
	Sku         string `json:"sku" example:"SKU123"`
	Description string `json:"description" example:"Cama 1 plaza"`
	Dimensions  struct {
		Length int64  `json:"length" example:"100"`
		Height int64  `json:"height" example:"100"`
		Unit   string `json:"unit" example:"cm"`
		Width  int64  `json:"width" example:"100"`
	} `json:"dimensions"`
	Weight struct {
		Unit  string `json:"unit" example:"g"`
		Value int64  `json:"value" example:"1800"`
	} `json:"weight"`
	Quantity struct {
		QuantityNumber int    `json:"quantityNumber" example:"1"`
		QuantityUnit   string `json:"quantityUnit" example:"unit"`
	} `json:"quantity"`
	Insurance struct {
		Currency  string `json:"currency" example:"CLP"`
		UnitValue int64  `json:"unitValue" example:"10000"`
	} `json:"insurance"`
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
