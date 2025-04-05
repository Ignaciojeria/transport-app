package mapper

import "transport-app/app/domain"

func MapItemReferencesToDomain(itemRefs []struct {
	Quantity struct {
		QuantityNumber int    `json:"quantityNumber"`
		QuantityUnit   string `json:"quantityUnit"`
	} `json:"quantity"`
	Sku string `json:"sku"`
}) []domain.ItemReference {
	mapped := make([]domain.ItemReference, len(itemRefs))
	for i, ref := range itemRefs {
		mapped[i] = domain.ItemReference{
			Sku: ref.Sku,
			Quantity: domain.Quantity{
				QuantityNumber: ref.Quantity.QuantityNumber,
				QuantityUnit:   ref.Quantity.QuantityUnit,
			},
		}
	}
	return mapped
}

func MapItemReferencesFromDomain(itemRefs []domain.ItemReference) []struct {
	Quantity struct {
		QuantityNumber int    `json:"quantityNumber"`
		QuantityUnit   string `json:"quantityUnit"`
	} `json:"quantity"`
	Sku string `json:"sku"`
} {
	mapped := make([]struct {
		Quantity struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		Sku string `json:"sku"`
	}, len(itemRefs))

	for i, ref := range itemRefs {
		mapped[i] = struct {
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			Sku string `json:"sku"`
		}{
			Quantity: struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			}{
				QuantityNumber: ref.Quantity.QuantityNumber,
				QuantityUnit:   ref.Quantity.QuantityUnit,
			},
			Sku: ref.Sku,
		}
	}

	return mapped
}
