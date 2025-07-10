package mapper

import (
	"transport-app/app/domain"

	"github.com/paulmach/orb"
)

func MapNodeInfoToDomain(nodeInfo struct {
	ReferenceID string `json:"referenceID"`
}, addressInfo struct {
	AddressLine1 string `json:"addressLine1" example:"Inglaterra 59"`
	AddressLine2 string `json:"addressLine2" example:"Piso 2214"`
	Contact      struct {
		AdditionalContactMethods []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"additionalContactMethods"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		NationalID string `json:"nationalID"`
		Documents  []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"documents"`
		FullName string `json:"fullName"`
	} `json:"contact"`
	Coordinates struct {
		Latitude   float64 `json:"latitude" example:"-33.5147889"`
		Longitude  float64 `json:"longitude" example:"-70.6130425"`
		Source     string  `json:"source" example:"GOOGLE_MAPS"`
		Confidence struct {
			Level   float64 `json:"level" example:"0.1"`
			Message string  `json:"message" example:"DISTRICT_CENTROID"`
			Reason  string  `json:"reason" example:"PROVIDER_RESULT_OUT_OF_DISTRICT"`
		} `json:"confidence"`
	} `json:"coordinates"`
	PoliticalArea struct {
		Code            string `json:"code" example:"cl-rm-la-florida"`
		AdminAreaLevel1 string `json:"adminAreaLevel1" example:"region metropolitana de santiago"`
		AdminAreaLevel2 string `json:"adminAreaLevel2" example:"santiago"`
		AdminAreaLevel3 string `json:"adminAreaLevel3" example:"la florida"`
		AdminAreaLevel4 string `json:"adminAreaLevel4" example:""`
		TimeZone        string `json:"timeZone" example:"America/Santiago"`
		Confidence      struct {
			Level   float64 `json:"level" example:"0.0"`
			Message string  `json:"message" example:""`
			Reason  string  `json:"reason" example:""`
		} `json:"confidence"`
	} `json:"politicalArea"`
	ZipCode string `json:"zipCode" example:"7500000"`
}) domain.NodeInfo {
	return domain.NodeInfo{
		ReferenceID: domain.ReferenceID(nodeInfo.ReferenceID),
		AddressInfo: domain.AddressInfo{
			Contact: domain.Contact{
				FullName:                 addressInfo.Contact.FullName,
				PrimaryEmail:             addressInfo.Contact.Email,
				PrimaryPhone:             addressInfo.Contact.Phone,
				NationalID:               addressInfo.Contact.NationalID,
				Documents:                MapDocumentsToDomain(addressInfo.Contact.Documents),
				AdditionalContactMethods: MapAdditionalContactMethodsToDomain(addressInfo.Contact.AdditionalContactMethods),
			},
			PoliticalArea: domain.PoliticalArea{
				Code:            addressInfo.PoliticalArea.Code,
				AdminAreaLevel1: addressInfo.PoliticalArea.AdminAreaLevel1,
				AdminAreaLevel2: addressInfo.PoliticalArea.AdminAreaLevel2,
				AdminAreaLevel3: addressInfo.PoliticalArea.AdminAreaLevel3,
				AdminAreaLevel4: addressInfo.PoliticalArea.AdminAreaLevel4,
				TimeZone:        addressInfo.PoliticalArea.TimeZone,
				Confidence: domain.CoordinatesConfidence{
					Level:   addressInfo.PoliticalArea.Confidence.Level,
					Message: addressInfo.PoliticalArea.Confidence.Message,
					Reason:  addressInfo.PoliticalArea.Confidence.Reason,
				},
			},
			AddressLine1: addressInfo.AddressLine1,
			AddressLine2: addressInfo.AddressLine2,
			Coordinates: domain.Coordinates{
				Point:  orb.Point{addressInfo.Coordinates.Longitude, addressInfo.Coordinates.Latitude},
				Source: addressInfo.Coordinates.Source,
				Confidence: domain.CoordinatesConfidence{
					Level:   addressInfo.Coordinates.Confidence.Level,
					Message: addressInfo.Coordinates.Confidence.Message,
					Reason:  addressInfo.Coordinates.Confidence.Reason,
				},
			},
			ZipCode: addressInfo.ZipCode,
		},
	}
}

func MapAdditionalContactMethodsToDomain(methods []struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}) []domain.ContactMethod {
	if methods == nil {
		return nil
	}

	domainMethods := make([]domain.ContactMethod, len(methods))
	for i, method := range methods {
		domainMethods[i] = domain.ContactMethod{
			Type:  method.Type,
			Value: method.Value,
		}
	}
	return domainMethods
}
