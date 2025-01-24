package domain

type Vehicle struct {
	Headers
	ID              int64
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
		ID           int64
		Organization Organization `json:"organization"`
		ReferenceID  string       `json:"referenceID"`
		Name         string       `json:"name"`
		NationalID   string       `json:"nationalID"`
	} `json:"carrier"`
}

func (v *Vehicle) UpdateIfChanged(in Vehicle) {
	if in.ReferenceID != "" {
		v.ReferenceID = in.ReferenceID
	}
	if in.Plate != "" {
		v.Plate = in.Plate
	}
	if in.CertificateDate != "" {
		v.CertificateDate = in.CertificateDate
	}
	if in.Category != "" {
		v.Category = in.Category
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
	if in.Carrier.ID != 0 {
		v.Carrier.ID = in.Carrier.ID
	}
	if in.Carrier.ReferenceID != "" {
		v.Carrier.ReferenceID = in.Carrier.ReferenceID
	}
	if in.Carrier.Name != "" {
		v.Carrier.Name = in.Carrier.Name
	}
	if in.Carrier.NationalID != "" {
		v.Carrier.NationalID = in.Carrier.NationalID
	}
}
