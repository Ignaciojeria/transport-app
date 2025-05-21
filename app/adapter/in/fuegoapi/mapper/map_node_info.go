package mapper

import (
	"transport-app/app/domain"

	"github.com/paulmach/orb"
)

func MapNodeInfoToDomain(nodeInfo struct {
	ReferenceID string `json:"referenceID"`
}, addressInfo struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
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

	District    string `json:"district"`
	Coordinates struct {
		Latitude             float64 `json:"latitude"`
		Longitude            float64 `json:"longitude"`
		Source               string  `json:"source"`
		RequiresManualReview bool    `json:"requiresManualReview"`
	} `json:"coordinates"`
	Province string `json:"province"`
	State    string `json:"state"`
	TimeZone string `json:"timeZone"`
	ZipCode  string `json:"zipCode"`
}) domain.NodeInfo {
	return domain.NodeInfo{
		ReferenceID: domain.ReferenceID(nodeInfo.ReferenceID),
		AddressInfo: domain.AddressInfo{
			Contact: domain.Contact{
				FullName:     addressInfo.Contact.FullName,
				PrimaryEmail: addressInfo.Contact.Email,
				PrimaryPhone: addressInfo.Contact.Phone,
				NationalID:   addressInfo.Contact.NationalID,
				Documents:    MapDocumentsToDomain(addressInfo.Contact.Documents),
			},
			State:        domain.State(addressInfo.State),
			Province:     domain.Province(addressInfo.Province),
			District:     domain.District(addressInfo.District),
			AddressLine1: addressInfo.AddressLine1,
			//AddressLine2: addressInfo.AddressLine2,
			//AddressLine3: addressInfo.AddressLine3,
			Location: orb.Point{
				addressInfo.Coordinates.Longitude, // orb.Point espera [lon, lat]
				addressInfo.Coordinates.Latitude,
			},
			RequiresManualReview: addressInfo.Coordinates.RequiresManualReview,
			CoordinateSource:     addressInfo.Coordinates.Source,
			ZipCode:              addressInfo.ZipCode,
			TimeZone:             addressInfo.TimeZone,
		},
	}
}

func MapNodeInfoToResponseNodeInfo(nodeInfo domain.NodeInfo) (struct {
	ReferenceID string `json:"referenceID"`
	Name        string `json:"name"`
}, struct {
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`

	Contact struct {
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		NationalID string `json:"nationalID"`
		Documents  []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"documents"`
		FullName string `json:"fullName"`
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

		Contact struct {
			Email      string `json:"email"`
			Phone      string `json:"phone"`
			NationalID string `json:"nationalID"`
			Documents  []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"documents"`
			FullName string `json:"fullName"`
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
		//	AddressLine2: nodeInfo.AddressInfo.AddressLine2,
		//	AddressLine3: nodeInfo.AddressInfo.AddressLine3,
		//	Locality:     nodeInfo.AddressInfo.Locality,
		District:  nodeInfo.AddressInfo.District.String(),
		Province:  nodeInfo.AddressInfo.Province.String(),
		State:     nodeInfo.AddressInfo.State.String(),
		ZipCode:   nodeInfo.AddressInfo.ZipCode,
		TimeZone:  nodeInfo.AddressInfo.TimeZone,
		Latitude:  nodeInfo.AddressInfo.Location[1], // Latitud
		Longitude: nodeInfo.AddressInfo.Location[0], // Longitud
	}

	// Mapear contacto
	responseAddressInfo.Contact.FullName = nodeInfo.AddressInfo.Contact.FullName
	responseAddressInfo.Contact.Email = nodeInfo.AddressInfo.Contact.PrimaryEmail
	responseAddressInfo.Contact.Phone = nodeInfo.AddressInfo.Contact.PrimaryPhone
	responseAddressInfo.Contact.NationalID = nodeInfo.AddressInfo.Contact.NationalID
	responseAddressInfo.Contact.Documents = MapDocumentsFromDomain(nodeInfo.AddressInfo.Contact.Documents)
	return responseNodeInfo, responseAddressInfo
}
