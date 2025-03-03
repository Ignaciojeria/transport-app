package request

import (
	"transport-app/app/domain"

	"github.com/paulmach/orb"
)

type UpsertNodeRequest struct {
	ReferenceID string `json:"referenceID" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	NodeAddress struct {
		AddressLine1 string  `json:"addressLine1"`
		AddressLine2 string  `json:"addressLine2"`
		AddressLine3 string  `json:"addressLine3"`
		Locality     string  `json:"locality"`
		District     string  `json:"district"`
		Latitude     float64 `json:"latitude"`
		Longitude    float64 `json:"longitude"`
		Province     string  `json:"province"`
		State        string  `json:"state"`
		TimeZone     string  `json:"timeZone"`
		ZipCode      string  `json:"zipCode"`
	} `json:"nodeAddress"`
	OperatorContact struct {
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		NationalID string `json:"nationalID"`
		Documents  []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"documents"`
		FullName string `json:"fullName"`
	} `json:"operatorContact"`
	References []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"references"`
}

func (req UpsertNodeRequest) Map() domain.NodeInfo {
	nodeInfo := domain.NodeInfo{
		ReferenceID: domain.ReferenceID(req.ReferenceID),
		Name:        req.Name,
		NodeType: domain.NodeType{
			Value: req.Type,
		},
		AddressInfo: domain.AddressInfo{
			AddressLine1: req.NodeAddress.AddressLine1,
			AddressLine2: req.NodeAddress.AddressLine2,
			AddressLine3: req.NodeAddress.AddressLine3,
			Locality:     req.NodeAddress.Locality,
			District:     req.NodeAddress.District,
			Location: orb.Point{
				req.NodeAddress.Longitude, // orb.Point espera [lon, lat]
				req.NodeAddress.Latitude,
			},
			Province: req.NodeAddress.Province,
			State:    req.NodeAddress.State,
			TimeZone: req.NodeAddress.TimeZone,
			ZipCode:  req.NodeAddress.ZipCode,
		},
		References: func() []domain.Reference {
			refs := make([]domain.Reference, len(req.References))
			for i, ref := range req.References {
				refs[i] = domain.Reference{
					Type:  ref.Type,
					Value: ref.Value,
				}
			}
			return refs
		}(),
		Contact: domain.Contact{
			Email:      req.OperatorContact.Email,
			Phone:      req.OperatorContact.Phone,
			NationalID: req.OperatorContact.NationalID,
			FullName:   req.OperatorContact.FullName,
			Documents: func() []domain.Document {
				docs := make([]domain.Document, len(req.OperatorContact.Documents))
				for i, doc := range req.OperatorContact.Documents {
					docs[i] = domain.Document{
						Type:  doc.Type,
						Value: doc.Value,
					}
				}
				return docs
			}(),
		},
	}
	return nodeInfo
}
