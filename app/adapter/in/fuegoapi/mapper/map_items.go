package mapper

import "transport-app/app/domain"

func MapItemsToDomain(items []struct {
	Description string `json:"description" example:"Cama 1 plaza"`
	Dimensions  struct {
		Length int64  `json:"length" example:"100"`
		Height int64  `json:"height" example:"100"`
		Unit   string `json:"unit" example:"cm"`
		Width  int64  `json:"width" example:"100"`
	} `json:"dimensions"`
	Insurance struct {
		Currency  string `json:"currency" example:"CLP"`
		UnitValue int64  `json:"unitValue" example:"10000"`
	} `json:"insurance"`
	Quantity struct {
		QuantityNumber int    `json:"quantityNumber" example:"1"`
		QuantityUnit   string `json:"quantityUnit" example:"unit"`
	} `json:"quantity"`
	Sku    string `json:"sku" example:"1234567890"`
	Weight struct {
		Unit  string `json:"unit" example:"g"`
		Value int64  `json:"value" example:"1800"`
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
