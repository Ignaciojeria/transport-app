package mapper

import "transport-app/app/domain"

// MapPackagesToDomain convierte estructuras anónimas de paquetes a domain.Package
func MapPackagesToDomain(packages []struct {
	Lpn       string `json:"lpn" example:"LPN456"`
	Volume    int64  `json:"volume" example:"1000"`
	Weight    int64  `json:"weight" example:"1000"`
	Insurance int64  `json:"insurance" example:"10000"`
	Skills    []string `json:"skills" example:"fragile"`
	Labels    []struct {
		Type  string `json:"type" example:"packageCode"`
		Value string `json:"value" example:"uuid"`
	} `json:"labels"`
	Items []struct {
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
	} `json:"items"`
}) []domain.DeliveryUnit {
	mapped := make([]domain.DeliveryUnit, len(packages))
	for i, pkg := range packages {
		// Calculate volume from items if not provided
		volume := pkg.Volume
		if volume == 0 && len(pkg.Items) > 0 {
			for _, item := range pkg.Items {
				itemVolume := item.Dimensions.Length * item.Dimensions.Width * item.Dimensions.Height
				volume += itemVolume * int64(item.Quantity.QuantityNumber)
			}
		}

		// Calculate weight from items if not provided
		weight := pkg.Weight
		if weight == 0 && len(pkg.Items) > 0 {
			for _, item := range pkg.Items {
				weight += item.Weight.Value * int64(item.Quantity.QuantityNumber)
			}
		}

		mapped[i] = domain.DeliveryUnit{
			Lpn:    pkg.Lpn,
			Volume: volume,
			Status: domain.Status{Status: domain.StatusAvailable},
			// Create dimensions from volume (assuming cubic root for simplicity)
			Dimensions: domain.Dimensions{
				Height: 0, // Will be calculated from items if needed
				Width:  0,
				Length: 0,
				Unit:   "cm",
			},
			Weight: domain.Weight{
				Unit:  "g", // Standard unit as per convention
				Value: weight,
			},
			Insurance: domain.Insurance{
				Currency:  "CLP", // Default currency, could be made configurable
				UnitValue: pkg.Insurance,
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
