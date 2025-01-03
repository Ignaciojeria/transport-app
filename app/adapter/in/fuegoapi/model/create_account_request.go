package model

import "transport-app/app/domain"

type CreateAccountRequest struct {
	NationalID string `json:"nationalID"`
	Email      string `json:"email"`
	Origin     struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
		} `json:"addressInfo"`
		NodeInfo struct {
			Name     string `json:"name"`
			Operator struct {
				Name        string `json:"name"`
				NationalID  string `json:"nationalId"`
				ReferenceID string `json:"referenceId"`
				Type        string `json:"type"`
			} `json:"operator"`
			ReferenceID string `json:"referenceId"`
			References  []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"references"`
			Type string `json:"type"`
		} `json:"nodeInfo"`
	} `json:"origin"`
}

func (req CreateAccountRequest) Map() domain.Account {
	return domain.Account{
		Contact: domain.Contact{
			FullName: req.Email, // Assuming email represents the user’s name or contact identifier
		},
		Origin: domain.Origin{
			AddressInfo: domain.AddressInfo{
				AddressLine1: req.Origin.AddressInfo.AddressLine1,
				AddressLine2: req.Origin.AddressInfo.AddressLine2,
			},
			NodeInfo: domain.NodeInfo{
				Name: req.Origin.NodeInfo.Name,
				Operator: domain.Operator{
					Name:        req.Origin.NodeInfo.Operator.Name,
					NationalID:  req.Origin.NodeInfo.Operator.NationalID,
					ReferenceID: domain.ReferenceID(req.Origin.NodeInfo.Operator.ReferenceID),
					Type:        req.Origin.NodeInfo.Operator.Type,
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
