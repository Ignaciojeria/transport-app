package request

type OutForDeliveryRequest struct {
	Plan struct {
		Routes []struct {
			Carrier struct {
				NationalID  string `json:"nationalID"`
				ReferenceID string `json:"referenceID"`
			} `json:"carrier"`
			Vehicle struct {
				Plate       string `json:"plate"`
				ReferenceID string `json:"referenceID"`
			} `json:"vehicle"`
			Orders []struct {
				ReferenceID         string `json:"referenceID"`
				BusinessIdentifiers struct {
					Commerce string `json:"commerce"`
					Consumer string `json:"consumer"`
				} `json:"businessIdentifiers"`
				OutForDeliveryDate string `json:"outForDeliveryDate"`
			} `json:"orders"`
			ReferenceID string `json:"referenceID"`
		} `json:"routes"`
	} `json:"plan"`
}
