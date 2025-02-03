package request

import "time"

type OrdersCheckoutRequest struct {
	OrderReferenceID string    `json:"orderReferenceID"`
	RouteReferenceID string    `json:"routeReferenceID"`
	Status           string    `json:"status"`
	DeliveredAt      time.Time `json:"deliveredAt"`
	Vehicle          struct {
		ReferenceID string `json:"referenceID"`
		Plate       string `json:"plate"`
		Carrier     struct {
			ReferenceID string `json:"referenceID"`
			NationalID  string `json:"nationalID"`
		} `json:"carrier"`
	} `json:"vehicle"`
	Recipient struct {
		FullName   string `json:"fullName"`
		NationalID string `json:"nationalID"`
	} `json:"recipient"`
	EvidencePhotos []struct {
		URL     string    `json:"url"`
		Type    string    `json:"type"`
		TakenAt time.Time `json:"takenAt"`
	} `json:"evidencePhotos"`
	Delivery struct {
		Failure struct {
			ReasonReferenceID string `json:"reasonReferenceID"`
			Reason            string `json:"reason"`
			Detail            string `json:"detail"`
		} `json:"failure"`
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
	} `json:"delivery"`
}
