package model

import "transport-app/app/domain"

type CreateAccountRequest struct {
	NationalID string `json:"nationalID" validate:"required" example:"18666636-4"`
	Email      string `json:"email" validate:"required" example:"ignaciovl.j@gmail.com"`
	Origin     struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1" validate:"required" example:"Inglaterra 59"`
			AddressLine2 string `json:"addressLine2" example:"La Florida, Región Metropolitana, Chile"`
		} `json:"addressInfo"`
		NodeInfo struct {
			Name     string `json:"name" example:"IGNACIO HUB"`
			Operator struct {
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
				Type string `json:"type" example:"ARRENDATARIO"`
			} `json:"operator"`
			ReferenceID string `json:"referenceId" example:"BODEGA_2214"`
			References  []struct {
				Type  string `json:"type" example:""`
				Value string `json:"value" example:""`
			} `json:"references"`
			Type string `json:"type" example:"BODEGA_DEPARTAMENTO"`
		} `json:"nodeInfo"`
	} `json:"origin"`
}

func (req CreateAccountRequest) Map() domain.Account {
	return domain.Account{
		NationalID: req.NationalID,
		Contact: domain.Contact{
			Email: req.Email,
		},
		Origin: domain.Origin{
			AddressInfo: domain.AddressInfo{
				AddressLine1: req.Origin.AddressInfo.AddressLine1,
				AddressLine2: req.Origin.AddressInfo.AddressLine2,
			},
			NodeInfo: domain.NodeInfo{
				Name: req.Origin.NodeInfo.Name,
				Operator: domain.Operator{
					Contact: domain.Contact{
						FullName: req.Origin.NodeInfo.Operator.Contact.FullName,
						Email:    req.Origin.NodeInfo.Operator.Contact.Email,
						Phone:    req.Origin.NodeInfo.Operator.Contact.Phone,
					},
					Type: req.Origin.NodeInfo.Operator.Type,
				},
				ReferenceID: domain.ReferenceID(req.Origin.NodeInfo.ReferenceID),
				References: func() []domain.References {
					references := []domain.References{}
					for _, ref := range req.Origin.NodeInfo.References {
						references = append(references, domain.References{
							Type:  ref.Type,
							Value: ref.Value,
						})
					}
					return references
				}(),
				Type: req.Origin.NodeInfo.Type,
			},
		},
		Profiles: []domain.Profile{}, // Assuming profiles will be populated later
	}
}
