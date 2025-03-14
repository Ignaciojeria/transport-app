package mapper

import (
	"transport-app/app/domain"

	"github.com/paulmach/orb"
)

func MapNodeInfoToDomain(nodeInfo struct {
	ReferenceID string `json:"referenceID"`
	Name        string `json:"name"`
}, addressInfo struct {
	ProviderAddress string `json:"providerAddress"`
	AddressLine1    string `json:"addressLine1"`
	AddressLine2    string `json:"addressLine2"`
	AddressLine3    string `json:"addressLine3"`
	Contact         struct {
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		NationalID string `json:"nationalID"`
		Documents  []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"documents"`
		FullName string `json:"fullName"`
	} `json:"contact"`
	Locality  string  `json:"locality"`
	District  string  `json:"district"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Province  string  `json:"province"`
	State     string  `json:"state"`
	TimeZone  string  `json:"timeZone"`
	ZipCode   string  `json:"zipCode"`
}) domain.NodeInfo {
	return domain.NodeInfo{
		ReferenceID: domain.ReferenceID(nodeInfo.ReferenceID),
		Name:        nodeInfo.Name,
		AddressInfo: domain.AddressInfo{
			Contact: domain.Contact{
				FullName:   addressInfo.Contact.FullName,
				Email:      addressInfo.Contact.Email,
				Phone:      addressInfo.Contact.Phone,
				NationalID: addressInfo.Contact.NationalID,
				Documents:  MapDocumentsToDomain(addressInfo.Contact.Documents),
			},
			State:        addressInfo.State,
			Locality:     addressInfo.Locality,
			Province:     addressInfo.Province,
			District:     addressInfo.District,
			AddressLine1: addressInfo.AddressLine1,
			AddressLine2: addressInfo.AddressLine2,
			AddressLine3: addressInfo.AddressLine3,
			Location: orb.Point{
				addressInfo.Longitude, // orb.Point espera [lon, lat]
				addressInfo.Latitude,
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
	ProviderAddress string `json:"providerAddress"`
	AddressLine1    string `json:"addressLine1"`
	AddressLine2    string `json:"addressLine2"`
	AddressLine3    string `json:"addressLine3"`
	Contact         struct {
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		NationalID string `json:"nationalID"`
		Documents  []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"documents"`
		FullName string `json:"fullName"`
	} `json:"contact"`
	Locality  string  `json:"locality"`
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
		ProviderAddress string `json:"providerAddress"`
		AddressLine1    string `json:"addressLine1"`
		AddressLine2    string `json:"addressLine2"`
		AddressLine3    string `json:"addressLine3"`
		Contact         struct {
			Email      string `json:"email"`
			Phone      string `json:"phone"`
			NationalID string `json:"nationalID"`
			Documents  []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"documents"`
			FullName string `json:"fullName"`
		} `json:"contact"`
		Locality  string  `json:"locality"`
		District  string  `json:"district"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Province  string  `json:"province"`
		State     string  `json:"state"`
		TimeZone  string  `json:"timeZone"`
		ZipCode   string  `json:"zipCode"`
	}{
		AddressLine1: nodeInfo.AddressInfo.AddressLine1,
		AddressLine2: nodeInfo.AddressInfo.AddressLine2,
		AddressLine3: nodeInfo.AddressInfo.AddressLine3,
		Locality:     nodeInfo.AddressInfo.Locality,
		District:     nodeInfo.AddressInfo.District,
		Province:     nodeInfo.AddressInfo.Province,
		State:        nodeInfo.AddressInfo.State,
		ZipCode:      nodeInfo.AddressInfo.ZipCode,
		TimeZone:     nodeInfo.AddressInfo.TimeZone,
		Latitude:     nodeInfo.AddressInfo.Location[1], // Latitud
		Longitude:    nodeInfo.AddressInfo.Location[0], // Longitud
	}

	// Mapear contacto
	responseAddressInfo.Contact.FullName = nodeInfo.Contact.FullName
	responseAddressInfo.Contact.Email = nodeInfo.Contact.Email
	responseAddressInfo.Contact.Phone = nodeInfo.Contact.Phone
	responseAddressInfo.Contact.NationalID = nodeInfo.Contact.NationalID
	responseAddressInfo.Contact.Documents = MapDocumentsFromDomain(nodeInfo.Contact.Documents)

	return responseNodeInfo, responseAddressInfo
}
