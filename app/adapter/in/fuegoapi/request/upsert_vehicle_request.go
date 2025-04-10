package request

import "transport-app/app/domain"

type UpsertVehicleRequest struct {
	ReferenceID     string `json:"referenceID" validate:"required"`
	Plate           string `json:"plate" validate:"required"`
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
		Name       string `json:"name"`
		NationalID string `json:"nationalID"`
	} `json:"carrier"`
}

func (v UpsertVehicleRequest) Map() domain.Vehicle {
	return domain.Vehicle{
		Plate:           v.Plate,
		IsActive:        v.IsActive,
		CertificateDate: v.CertificateDate,
		VehicleCategory: domain.VehicleCategory{
			Type: v.Category,
		},
		Weight: struct {
			Value         int
			UnitOfMeasure string
		}{
			Value:         v.Weight.Value,
			UnitOfMeasure: v.Weight.UnitOfMeasure,
		},
		Insurance: struct {
			PolicyStartDate      string
			PolicyExpirationDate string
			PolicyRenewalDate    string
			MaxInsuranceCoverage struct {
				Amount   float64
				Currency string
			}
		}{
			PolicyStartDate:      v.Insurance.PolicyStartDate,
			PolicyExpirationDate: v.Insurance.PolicyExpirationDate,
			PolicyRenewalDate:    v.Insurance.PolicyRenewalDate,
			MaxInsuranceCoverage: struct {
				Amount   float64
				Currency string
			}{
				Amount:   v.Insurance.MaxInsuranceCoverage.Amount,
				Currency: v.Insurance.MaxInsuranceCoverage.Currency,
			},
		},
		TechnicalReview: struct {
			LastReviewDate string
			NextReviewDate string
			ReviewedBy     string
		}{
			LastReviewDate: v.TechnicalReview.LastReviewDate,
			NextReviewDate: v.TechnicalReview.NextReviewDate,
			ReviewedBy:     v.TechnicalReview.ReviewedBy,
		},
		Dimensions: struct {
			Width         float64
			Length        float64
			Height        int
			UnitOfMeasure string
		}{
			Width:         v.Dimensions.Width,
			Length:        v.Dimensions.Length,
			Height:        v.Dimensions.Height,
			UnitOfMeasure: v.Dimensions.UnitOfMeasure,
		},
		Carrier: struct {
			Name       string
			NationalID string
		}{
			Name:       v.Carrier.Name,
			NationalID: v.Carrier.NationalID,
		},
	}
}
