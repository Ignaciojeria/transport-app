package domain

type Vehicle struct {
	Headers
	Plate           string          `json:"plate"`
	IsActive        bool            `json:"isActive"`
	CertificateDate string          `json:"certificateDate"`
	VehicleCategory VehicleCategory `json:"category"`
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
	Carrier Carrier `json:"carrier"`
}

func (v Vehicle) UpdateIfChanged(in Vehicle) Vehicle {
	v.Organization = in.Organization

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

type VehicleCategory struct {
	Organization        Organization
	ID                  int64
	Type                string
	MaxPackagesQuantity int
}

func (vc VehicleCategory) UpdateIfChanged(in VehicleCategory) VehicleCategory {
	if in.Type != "" {
		vc.Type = in.Type
	}
	if in.MaxPackagesQuantity != 0 {
		vc.MaxPackagesQuantity = in.MaxPackagesQuantity
	}
	if in.Organization.ID != 0 {
		vc.Organization = in.Organization
	}
	if in.ID != 0 {
		vc.ID = in.ID
	}
	return vc
}
