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

	District    string `json:"district" example:"la florida"`
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
	Province string `json:"province" example:"santiago"`
	State    string `json:"state" example:"region metropolitana de santiago"`
	TimeZone string `json:"timeZone" example:"America/Santiago"`
	ZipCode  string `json:"zipCode" example:"7500000"`
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
			State:        domain.State(addressInfo.State),
			Province:     domain.Province(addressInfo.Province),
			District:     domain.District(addressInfo.District),
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
			ZipCode:  addressInfo.ZipCode,
			TimeZone: addressInfo.TimeZone,
		},
	}
}

func MapNodeInfoToResponseNodeInfo(nodeInfo domain.NodeInfo) (struct {
	ReferenceID string `json:"referenceID"`
	Name        string `json:"name"`
}, struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	Contact      struct {
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		NationalID string `json:"nationalID"`
		Documents  []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"documents"`
		FullName                 string `json:"fullName"`
		AdditionalContactMethods []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"additionalContactMethods"`
	} `json:"contact"`
	District  string  `json:"district"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Province  string  `json:"province"`
	State     string  `json:"state"`
	TimeZone  string  `json:"timeZone"`
	ZipCode   string  `json:"zipCode"`
}) {
	responseNodeInfo := struct {
		ReferenceID string `json:"referenceID"`
		Name        string `json:"name"`
	}{
		ReferenceID: string(nodeInfo.ReferenceID),
		Name:        nodeInfo.Name,
	}

	responseAddressInfo := struct {
		AddressLine1 string `json:"addressLine1"`
		AddressLine2 string `json:"addressLine2"`
		Contact      struct {
			Email      string `json:"email"`
			Phone      string `json:"phone"`
			NationalID string `json:"nationalID"`
			Documents  []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"documents"`
			FullName                 string `json:"fullName"`
			AdditionalContactMethods []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"additionalContactMethods"`
		} `json:"contact"`
		District  string  `json:"district"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Province  string  `json:"province"`
		State     string  `json:"state"`
		TimeZone  string  `json:"timeZone"`
		ZipCode   string  `json:"zipCode"`
	}{
		AddressLine1: nodeInfo.AddressInfo.AddressLine1,
		District:     nodeInfo.AddressInfo.District.String(),
		Province:     nodeInfo.AddressInfo.Province.String(),
		State:        nodeInfo.AddressInfo.State.String(),
		ZipCode:      nodeInfo.AddressInfo.ZipCode,
		TimeZone:     nodeInfo.AddressInfo.TimeZone,
		Latitude:     nodeInfo.AddressInfo.Coordinates.Point.Lat(),
		Longitude:    nodeInfo.AddressInfo.Coordinates.Point.Lon(),
	}

	responseAddressInfo.Contact.FullName = nodeInfo.AddressInfo.Contact.FullName
	responseAddressInfo.Contact.Email = nodeInfo.AddressInfo.Contact.PrimaryEmail
	responseAddressInfo.Contact.Phone = nodeInfo.AddressInfo.Contact.PrimaryPhone
	responseAddressInfo.Contact.NationalID = nodeInfo.AddressInfo.Contact.NationalID
	responseAddressInfo.Contact.Documents = MapDocumentsFromDomain(nodeInfo.AddressInfo.Contact.Documents)
	responseAddressInfo.Contact.AdditionalContactMethods = MapAdditionalContactMethodsFromDomain(nodeInfo.AddressInfo.Contact.AdditionalContactMethods)
	return responseNodeInfo, responseAddressInfo
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

func MapAdditionalContactMethodsFromDomain(methods []domain.ContactMethod) []struct {
	Type  string `json:"type"`
	Value string `json:"value"`
} {
	if methods == nil {
		return nil
	}

	responseMethods := make([]struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}, len(methods))

	for i, method := range methods {
		responseMethods[i] = struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}{
			Type:  method.Type,
			Value: method.Value,
		}
	}
	return responseMethods
}
