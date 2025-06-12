package request

import (
	"context"
	"transport-app/app/domain"
)

type OrderDestinationFixRequest struct {
	ManualChange struct {
		PerformedBy string `json:"performedBy" example:"juan@example.com"`
		Reason      string `json:"reason" example:"PROVIDER_RESULT_OUT_OF_DISTRICT"`
	} `json:"manualChange"`
	Destination struct {
		AddressLine1 string `json:"addressLine1" example:"Inglaterra 59"`
		AddressLine2 string `json:"addressLine2" example:"Piso 2214"`
		Coordinates  struct {
			Confidence struct {
				Level   float64 `json:"level" example:"0.1"`
				Message string  `json:"message" example:"DISTRICT_CENTROID"`
				Reason  string  `json:"reason" example:"PROVIDER_RESULT_OUT_OF_DISTRICT"`
			} `json:"confidence"`
			Latitude  float64 `json:"latitude" example:"-33.5147889"`
			Longitude float64 `json:"longitude" example:"-70.6130425"`
			Source    string  `json:"source" example:"GOOGLE_MAPS"`
		} `json:"coordinates"`
		ZipCode       string `json:"zipCode" example:"7500000"`
		PoliticalArea struct {
			Code       string `json:"id" example:"cl-rm-la-florida"`
			Province   string `json:"province" example:"santiago"`
			State      string `json:"state" example:"region metropolitana de santiago"`
			District   string `json:"district" example:"la florida"`
			TimeZone   string `json:"timeZone" example:"America/Santiago"`
			Confidence struct {
				Level   float64 `json:"level" example:"0.0"`
				Message string  `json:"message" example:""`
				Reason  string  `json:"reason" example:""`
			} `json:"confidence"`
		} `json:"politicalArea"`
	} `json:"destination"`
	OrderReferenceIDs []struct {
		BusinessIdentifiers struct {
			Commerce string `json:"commerce" example:"string"`
			Consumer string `json:"consumer" example:"string"`
		} `json:"businessIdentifiers"`
		ReferenceID string `json:"referenceId" example:"12345"`
	} `json:"orderReferenceIDs"`
}

// Map convierte el request a un objeto de dominio OrderDestinationFix
func (req OrderDestinationFixRequest) Map(ctx context.Context) []domain.Order {
	return nil
}
