package mapper

import "transport-app/app/domain"

// MapPackagesToDomain convierte estructuras anónimas de paquetes a domain.Package
func MapPackagesToDomain(packages []struct {
	SizeCategory string `json:"sizeCategory" example:"XL"`
	Dimensions   struct {
		Length int64  `json:"length" example:"100"`
		Height int64  `json:"height" example:"100"`
		Unit   string `json:"unit" example:"cm"`
		Width  int64  `json:"width" example:"100"`
	} `json:"dimensions"`
	Insurance struct {
		Currency  string `json:"currency" example:"CLP"`
		UnitValue int64  `json:"unitValue" example:"10000"`
	} `json:"insurance"`
	Items []struct {
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
	} `json:"items"`
	Lpn    string `json:"lpn" example:"1234567890"`
	Labels []struct {
		Type  string `json:"type" example:"packageCode"`
		Value string `json:"value" example:"uuid"`
	} `json:"labels"`
	Skills []string `json:"skills"`
	Weight struct {
		Unit  string `json:"unit" example:"g"`
		Value int64  `json:"value" example:"1800"`
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
			Skills: MapSkillsToDomain(pkg.Skills),
		}
	}
	return mapped
}

func MapSkillsToDomain(skills []string) []domain.Skill {
	mapped := make([]domain.Skill, len(skills))
	for i, skill := range skills {
		mapped[i] = domain.Skill(skill)
	}
	return mapped
}

func MapLabelsToDomain(labels []struct {
	Type  string `json:"type" example:"packageCode"`
	Value string `json:"value" example:"uuid"`
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
	Type  string `json:"type" example:"packageCode"`
	Value string `json:"value" example:"uuid"`
} {
	mapped := make([]struct {
		Type  string `json:"type" example:"packageCode"`
		Value string `json:"value" example:"uuid"`
	}, len(labels))
	for i, label := range labels {
		mapped[i] = struct {
			Type  string `json:"type" example:"packageCode"`
			Value string `json:"value" example:"uuid"`
		}{
			Type:  label.Type,
			Value: label.Value,
		}
	}
	return mapped
}
