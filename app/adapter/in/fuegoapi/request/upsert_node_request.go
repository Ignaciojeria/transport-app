package request

import (
	"transport-app/app/domain"

	"github.com/paulmach/orb"
)

type UpsertNodeRequest struct {
	ReferenceID string `json:"referenceID"`
	Name        string `json:"name"`
	NodeAddress struct {
		AddressLine1 string `json:"addressLine1"`
		AddressLine2 string `json:"addressLine2"`
		Coordinates  struct {
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
			Code            string `json:"id" example:"cl-rm-la-florida"`
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
		ZipCode string `json:"zipCode"`
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
			PoliticalArea: domain.PoliticalArea{
				Code:            r.NodeAddress.PoliticalArea.Code,
				AdminAreaLevel1: r.NodeAddress.PoliticalArea.AdminAreaLevel1,
				AdminAreaLevel2: r.NodeAddress.PoliticalArea.AdminAreaLevel2,
				AdminAreaLevel3: r.NodeAddress.PoliticalArea.AdminAreaLevel3,
				AdminAreaLevel4: r.NodeAddress.PoliticalArea.AdminAreaLevel4,
				TimeZone:        r.NodeAddress.PoliticalArea.TimeZone,
			},
			AddressLine1: r.NodeAddress.AddressLine1,
			AddressLine2: r.NodeAddress.AddressLine2,
			Coordinates: domain.Coordinates{
				Point: orb.Point{
					r.NodeAddress.Coordinates.Longitude,
					r.NodeAddress.Coordinates.Latitude,
				},
				Source: r.NodeAddress.Coordinates.Source,
				Confidence: domain.CoordinatesConfidence{
					Level:   r.NodeAddress.Coordinates.Confidence.Level,
					Message: r.NodeAddress.Coordinates.Confidence.Message,
					Reason:  r.NodeAddress.Coordinates.Confidence.Reason,
				},
			},
			ZipCode: r.NodeAddress.ZipCode,
		},
	}
}
