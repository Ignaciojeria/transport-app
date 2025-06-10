package request

import (
	"transport-app/app/domain"

	"github.com/paulmach/orb"
)

type UpsertNodeRequest struct {
	ReferenceID string `json:"referenceID"`
	Name        string `json:"name"`
	NodeAddress struct {
		AddressLine1 string  `json:"addressLine1"`
		AddressLine2 string  `json:"addressLine2"`
		District     string  `json:"district"`
		Latitude     float64 `json:"latitude"`
		Locality     string  `json:"locality"`
		Longitude    float64 `json:"longitude"`
		Province     string  `json:"province"`
		State        string  `json:"state"`
		TimeZone     string  `json:"timeZone"`
		ZipCode      string  `json:"zipCode"`
	} `json:"nodeAddress"`
	Contact struct {
		Documents []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"documents"`
		Email      string `json:"email"`
		FullName   string `json:"fullName"`
		NationalID string `json:"nationalID"`
		Phone      string `json:"phone"`
	} `json:"contact"`
	References []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"references"`
	Type string `json:"type"`
}

func (r UpsertNodeRequest) Map() domain.NodeInfo {
	// Convertir documentos
	documents := make([]domain.Document, len(r.Contact.Documents))
	for i, doc := range r.Contact.Documents {
		documents[i] = domain.Document{
			Type:  doc.Type,
			Value: doc.Value,
		}
	}

	// Convertir referencias
	references := make([]domain.Reference, len(r.References))
	for i, ref := range r.References {
		references[i] = domain.Reference{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}

	// Crear el objeto NodeInfo
	return domain.NodeInfo{
		ReferenceID: domain.ReferenceID(r.ReferenceID),
		Name:        r.Name,
		NodeType: domain.NodeType{
			Value: r.Type,
		},
		References: references,
		AddressInfo: domain.AddressInfo{
			Contact: domain.Contact{
				Documents:    documents,
				FullName:     r.Contact.FullName,
				PrimaryEmail: r.Contact.Email,
				PrimaryPhone: r.Contact.Phone,
				NationalID:   r.Contact.NationalID,
			},
			State:        domain.State(r.NodeAddress.State),
			Province:     domain.Province(r.NodeAddress.Province),
			District:     domain.District(r.NodeAddress.District),
			AddressLine1: r.NodeAddress.AddressLine1,
			AddressLine2: r.NodeAddress.AddressLine2,
			Coordinates: domain.Coordinates{
				Point: orb.Point{r.NodeAddress.Longitude, r.NodeAddress.Latitude},
			},
			ZipCode:  r.NodeAddress.ZipCode,
			TimeZone: r.NodeAddress.TimeZone,
		},
	}
}
