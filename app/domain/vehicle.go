package domain

type Vehicle struct {
	ID              int64
	Organization    Organization `json:"organization"`
	ReferenceID     string       `json:"referenceID"`
	Plate           string       `json:"plate"`
	IsActive        bool         `json:"isActive"`
	CertificateDate string       `json:"certificateDate"`
	Category        string       `json:"category"`
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
		Document     struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"document"`
	} `json:"carrier"`
}
