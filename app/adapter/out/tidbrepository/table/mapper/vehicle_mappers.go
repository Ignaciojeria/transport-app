package mapper

import (
	"encoding/json"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func DomainToTableVehicle(d domain.Vehicle) table.Vehicle {
	weight, _ := json.Marshal(d.Weight)
	insurance, _ := json.Marshal(d.Insurance)
	technicalReview, _ := json.Marshal(d.TechnicalReview)
	dimensions, _ := json.Marshal(d.Dimensions)

	return table.Vehicle{
		ID:              d.ID,
		ReferenceID:     d.ReferenceID,
		Plate:           d.Plate,
		IsActive:        d.IsActive,
		CertificateDate: d.CertificateDate,
		Category:        d.Category,
		Weight:          table.JSONB(weight),
		Insurance:       table.JSONB(insurance),
		TechnicalReview: table.JSONB(technicalReview),
		Dimensions:      table.JSONB(dimensions),
		CarrierID:       d.Carrier.ID,
		Carrier: table.Carrier{
			ID:          d.Carrier.ID,
			ReferenceID: d.Carrier.ReferenceID,
			Name:        d.Carrier.Name,
			NationalID:  d.Carrier.NationalID,
		},
	}
}

func TableToDomainVehicle(t table.Vehicle, carrier table.Carrier) domain.Vehicle {
	var weight struct {
		Value         int    `json:"value"`
		UnitOfMeasure string `json:"unitOfMeasure"`
	}
	var insurance struct {
		PolicyStartDate      string `json:"policyStartDate"`
		PolicyExpirationDate string `json:"policyExpirationDate"`
		PolicyRenewalDate    string `json:"policyRenewalDate"`
		MaxInsuranceCoverage struct {
			Amount   float64 `json:"amount"`
			Currency string  `json:"currency"`
		} `json:"maxInsuranceCoverage"`
	}
	var technicalReview struct {
		LastReviewDate string `json:"lastReviewDate"`
		NextReviewDate string `json:"nextReviewDate"`
		ReviewedBy     string `json:"reviewedBy"`
	}
	var dimensions struct {
		Width         float64 `json:"width"`
		Length        float64 `json:"length"`
		Height        int     `json:"height"`
		UnitOfMeasure string  `json:"unitOfMeasure"`
	}
	var document struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	// Deserializar JSON de las tablas
	_ = json.Unmarshal(t.Weight, &weight)
	_ = json.Unmarshal(t.Insurance, &insurance)
	_ = json.Unmarshal(t.TechnicalReview, &technicalReview)
	_ = json.Unmarshal(t.Dimensions, &dimensions)
	_ = json.Unmarshal(carrier.Document, &document)
	return domain.Vehicle{
		ID:              t.ID,
		ReferenceID:     t.ReferenceID,
		Plate:           t.Plate,
		IsActive:        t.IsActive,
		CertificateDate: t.CertificateDate,
		Category:        t.Category,
		Weight:          weight,
		Insurance:       insurance,
		TechnicalReview: technicalReview,
		Dimensions:      dimensions,
		Carrier: struct {
			ID           int64
			Organization domain.Organization `json:"organization"`
			ReferenceID  string              `json:"referenceID"`
			Name         string              `json:"name"`
			NationalID   string              `json:"nationalID"`
		}{
			ID:           carrier.ID,
			Organization: domain.Organization{}, // Ajustar según tu lógica para inicializar Organization
			ReferenceID:  t.Carrier.ReferenceID,
			Name:         carrier.Name,
			NationalID:   carrier.NationalID,
		},
	}
}
