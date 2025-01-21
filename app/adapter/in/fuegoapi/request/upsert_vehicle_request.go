package request

import "transport-app/app/domain"

type UpsertVehicleRequest struct {
	ReferenceID     string `json:"referenceID"`
	Plate           string `json:"plate"`
	IsActive        bool   `json:"isActive"`
	CertificateDate string `json:"certificateDate"`
	Category        string `json:"category"`
	Weight          struct {
		Value         int    `json:"value"`
		UnitOfMeasure string `json:"unitOfMeasure"`
	} `json:"weight"`
	Insurance struct {
		PolicyStartDate      string `json:"policyStartDate"`
		PolicyExpirationDate string `json:"policyExpirationDate"`
		PolicyRenewalDate    string `json:"policyRenewalDate"`
		MaxInsuranceCoverage struct {
			Amount   float64 `json:"amount"`
			Currency string  `json:"currency"`
		} `json:"maxInsuranceCoverage"`
	} `json:"insurance"`
	TechnicalReview struct {
		LastReviewDate string `json:"lastReviewDate"`
		NextReviewDate string `json:"nextReviewDate"`
		ReviewedBy     string `json:"reviewedBy"`
	} `json:"technicalReview"`
	Dimensions struct {
		Width         float64 `json:"width"`
		Length        float64 `json:"length"`
		Height        int     `json:"height"`
		UnitOfMeasure string  `json:"unitOfMeasure"`
	} `json:"dimensions"`
	Carrier struct {
		ReferenceID string `json:"referenceID"`
		Name        string `json:"name"`
		NationalID  string `json:"nationalID"`
		Document    struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"document"`
	} `json:"carrier"`
}

func (v UpsertVehicleRequest) Map() domain.Vehicle {
	return domain.Vehicle{
		ReferenceID:     v.ReferenceID,
		Plate:           v.Plate,
		IsActive:        v.IsActive,
		CertificateDate: v.CertificateDate,
		Category:        v.Category,
		Weight: struct {
			Value         int    `json:"value"`
			UnitOfMeasure string `json:"unitOfMeasure"`
		}{
			Value:         v.Weight.Value,
			UnitOfMeasure: v.Weight.UnitOfMeasure,
		},
		Insurance: struct {
			PolicyStartDate      string `json:"policyStartDate"`
			PolicyExpirationDate string `json:"policyExpirationDate"`
			PolicyRenewalDate    string `json:"policyRenewalDate"`
			MaxInsuranceCoverage struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
			} `json:"maxInsuranceCoverage"`
		}{
			PolicyStartDate:      v.Insurance.PolicyStartDate,
			PolicyExpirationDate: v.Insurance.PolicyExpirationDate,
			PolicyRenewalDate:    v.Insurance.PolicyRenewalDate,
			MaxInsuranceCoverage: struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
			}{
				Amount:   v.Insurance.MaxInsuranceCoverage.Amount,
				Currency: v.Insurance.MaxInsuranceCoverage.Currency,
			},
		},
		TechnicalReview: struct {
			LastReviewDate string `json:"lastReviewDate"`
			NextReviewDate string `json:"nextReviewDate"`
			ReviewedBy     string `json:"reviewedBy"`
		}{
			LastReviewDate: v.TechnicalReview.LastReviewDate,
			NextReviewDate: v.TechnicalReview.NextReviewDate,
			ReviewedBy:     v.TechnicalReview.ReviewedBy,
		},
		Dimensions: struct {
			Width         float64 `json:"width"`
			Length        float64 `json:"length"`
			Height        int     `json:"height"`
			UnitOfMeasure string  `json:"unitOfMeasure"`
		}{
			Width:         v.Dimensions.Width,
			Length:        v.Dimensions.Length,
			Height:        v.Dimensions.Height,
			UnitOfMeasure: v.Dimensions.UnitOfMeasure,
		},
		Carrier: struct {
			ID           int64
			Organization domain.Organization `json:"organization"`
			ReferenceID  string              `json:"referenceID"`
			Name         string              `json:"name"`
			NationalID   string              `json:"nationalID"`
		}{
			ID:           0,                     // Si el ID no está presente en la solicitud
			Organization: domain.Organization{}, // Ajustar según tu lógica de negocio
			ReferenceID:  v.Carrier.ReferenceID,
			Name:         v.Carrier.Name,
			NationalID:   v.Carrier.NationalID,
		},
		Organization: domain.Organization{}, // Ajustar según corresponda
	}
}
