package domain

type Vehicle struct {
	Headers
	Plate           string
	IsActive        bool
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

func (v Vehicle) DocID() DocumentID {
	//return Hash(v.Organization, v.Plate)
	return ""
}

func (v Vehicle) UpdateIfChanged(in Vehicle) Vehicle {
	if in.Plate != "" {
		v.Plate = in.Plate
	}
	if in.CertificateDate != "" {
		v.CertificateDate = in.CertificateDate
	}
	if in.Weight.Value != 0 {
		v.Weight.Value = in.Weight.Value
	}
	if in.Weight.UnitOfMeasure != "" {
		v.Weight.UnitOfMeasure = in.Weight.UnitOfMeasure
	}

	if in.Insurance.PolicyStartDate != "" {
		v.Insurance.PolicyStartDate = in.Insurance.PolicyStartDate
	}
	if in.Insurance.PolicyExpirationDate != "" {
		v.Insurance.PolicyExpirationDate = in.Insurance.PolicyExpirationDate
	}
	if in.Insurance.PolicyRenewalDate != "" {
		v.Insurance.PolicyRenewalDate = in.Insurance.PolicyRenewalDate
	}
	if in.Insurance.MaxInsuranceCoverage.Amount != 0 {
		v.Insurance.MaxInsuranceCoverage.Amount = in.Insurance.MaxInsuranceCoverage.Amount
	}
	if in.Insurance.MaxInsuranceCoverage.Currency != "" {
		v.Insurance.MaxInsuranceCoverage.Currency = in.Insurance.MaxInsuranceCoverage.Currency
	}

	if in.TechnicalReview.LastReviewDate != "" {
		v.TechnicalReview.LastReviewDate = in.TechnicalReview.LastReviewDate
	}
	if in.TechnicalReview.NextReviewDate != "" {
		v.TechnicalReview.NextReviewDate = in.TechnicalReview.NextReviewDate
	}
	if in.TechnicalReview.ReviewedBy != "" {
		v.TechnicalReview.ReviewedBy = in.TechnicalReview.ReviewedBy
	}

	if in.Dimensions.Width != 0 {
		v.Dimensions.Width = in.Dimensions.Width
	}
	if in.Dimensions.Length != 0 {
		v.Dimensions.Length = in.Dimensions.Length
	}
	if in.Dimensions.Height != 0 {
		v.Dimensions.Height = in.Dimensions.Height
	}
	if in.Dimensions.UnitOfMeasure != "" {
		v.Dimensions.UnitOfMeasure = in.Dimensions.UnitOfMeasure
	}
	if in.VehicleCategory.Type != "" {
		v.VehicleCategory = in.VehicleCategory
	}
	return v
}
