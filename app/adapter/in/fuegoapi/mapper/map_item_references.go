package mapper

import "transport-app/app/domain"

func MapItemReferencesToDomain(itemRefs []struct {
	Quantity struct {
		QuantityNumber int    `json:"quantityNumber"`
		QuantityUnit   string `json:"quantityUnit"`
	} `json:"quantity"`
	ReferenceID string `json:"referenceId"`
}) []domain.ItemReference {
	mapped := make([]domain.ItemReference, len(itemRefs))
	for i, ref := range itemRefs {
		mapped[i] = domain.ItemReference{
			ReferenceID: domain.ReferenceID(ref.ReferenceID),
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
	ReferenceID string `json:"referenceId"`
} {
	mapped := make([]struct {
		Quantity struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		ReferenceID string `json:"referenceId"`
	}, len(itemRefs))

	for i, ref := range itemRefs {
		mapped[i] = struct {
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			ReferenceID string `json:"referenceId"`
		}{
			Quantity: struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			}{
				QuantityNumber: ref.Quantity.QuantityNumber,
				QuantityUnit:   ref.Quantity.QuantityUnit,
			},
			ReferenceID: string(ref.ReferenceID),
		}
	}

	return mapped
}
