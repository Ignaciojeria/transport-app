package domain

import "context"

type Vehicle struct {
	Headers
	Plate           string
	CertificateDate string
	VehicleCategory VehicleCategory
	Weight          struct {
		Value         int
		UnitOfMeasure string
	}
	Insurance struct {
		PolicyStartDate      string
		PolicyExpirationDate string
		PolicyRenewalDate    string
		MaxInsuranceCoverage struct {
			Amount   float64
			Currency string
		}
	}
	TechnicalReview struct {
		LastReviewDate string
		NextReviewDate string
		ReviewedBy     string
	}
	Dimensions struct {
		Width         float64
		Length        float64
		Height        int
		UnitOfMeasure string
	}
	Carrier Carrier
}

func (v Vehicle) DocID(ctx context.Context) DocumentID {
	return Hash(ctx, v.Plate)
}

func (v Vehicle) UpdateIfChanged(in Vehicle) (Vehicle, bool) {
	changed := false

	if in.Plate != "" && in.Plate != v.Plate {
		v.Plate = in.Plate
		changed = true
	}
	if in.CertificateDate != "" && in.CertificateDate != v.CertificateDate {
		v.CertificateDate = in.CertificateDate
		changed = true
	}
	if in.Weight.Value != 0 && in.Weight.Value != v.Weight.Value {
		v.Weight.Value = in.Weight.Value
		changed = true
	}
	if in.Weight.UnitOfMeasure != "" && in.Weight.UnitOfMeasure != v.Weight.UnitOfMeasure {
		v.Weight.UnitOfMeasure = in.Weight.UnitOfMeasure
		changed = true
	}

	if in.Insurance.PolicyStartDate != "" && in.Insurance.PolicyStartDate != v.Insurance.PolicyStartDate {
		v.Insurance.PolicyStartDate = in.Insurance.PolicyStartDate
		changed = true
	}
	if in.Insurance.PolicyExpirationDate != "" && in.Insurance.PolicyExpirationDate != v.Insurance.PolicyExpirationDate {
		v.Insurance.PolicyExpirationDate = in.Insurance.PolicyExpirationDate
		changed = true
	}
	if in.Insurance.PolicyRenewalDate != "" && in.Insurance.PolicyRenewalDate != v.Insurance.PolicyRenewalDate {
		v.Insurance.PolicyRenewalDate = in.Insurance.PolicyRenewalDate
		changed = true
	}
	if in.Insurance.MaxInsuranceCoverage.Amount != 0 && in.Insurance.MaxInsuranceCoverage.Amount != v.Insurance.MaxInsuranceCoverage.Amount {
		v.Insurance.MaxInsuranceCoverage.Amount = in.Insurance.MaxInsuranceCoverage.Amount
		changed = true
	}
	if in.Insurance.MaxInsuranceCoverage.Currency != "" && in.Insurance.MaxInsuranceCoverage.Currency != v.Insurance.MaxInsuranceCoverage.Currency {
		v.Insurance.MaxInsuranceCoverage.Currency = in.Insurance.MaxInsuranceCoverage.Currency
		changed = true
	}

	if in.TechnicalReview.LastReviewDate != "" && in.TechnicalReview.LastReviewDate != v.TechnicalReview.LastReviewDate {
		v.TechnicalReview.LastReviewDate = in.TechnicalReview.LastReviewDate
		changed = true
	}
	if in.TechnicalReview.NextReviewDate != "" && in.TechnicalReview.NextReviewDate != v.TechnicalReview.NextReviewDate {
		v.TechnicalReview.NextReviewDate = in.TechnicalReview.NextReviewDate
		changed = true
	}
	if in.TechnicalReview.ReviewedBy != "" && in.TechnicalReview.ReviewedBy != v.TechnicalReview.ReviewedBy {
		v.TechnicalReview.ReviewedBy = in.TechnicalReview.ReviewedBy
		changed = true
	}

	if in.Dimensions.Width != 0 && in.Dimensions.Width != v.Dimensions.Width {
		v.Dimensions.Width = in.Dimensions.Width
		changed = true
	}
	if in.Dimensions.Length != 0 && in.Dimensions.Length != v.Dimensions.Length {
		v.Dimensions.Length = in.Dimensions.Length
		changed = true
	}
	if in.Dimensions.Height != 0 && in.Dimensions.Height != v.Dimensions.Height {
		v.Dimensions.Height = in.Dimensions.Height
		changed = true
	}
	if in.Dimensions.UnitOfMeasure != "" && in.Dimensions.UnitOfMeasure != v.Dimensions.UnitOfMeasure {
		v.Dimensions.UnitOfMeasure = in.Dimensions.UnitOfMeasure
		changed = true
	}

	// Para VehicleCategory, actualizamos cada campo individualmente
	if in.VehicleCategory.Type != "" && in.VehicleCategory.Type != v.VehicleCategory.Type {
		v.VehicleCategory.Type = in.VehicleCategory.Type
		changed = true
	}
	if in.VehicleCategory.MaxPackagesQuantity != 0 && in.VehicleCategory.MaxPackagesQuantity != v.VehicleCategory.MaxPackagesQuantity {
		v.VehicleCategory.MaxPackagesQuantity = in.VehicleCategory.MaxPackagesQuantity
		changed = true
	}

	return v, changed
}
