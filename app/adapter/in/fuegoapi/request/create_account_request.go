package request

import "transport-app/app/domain"

type CreateAccountRequest struct {
	Contact struct {
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		NationalID string `json:"nationalID"`
		FullName   string `json:"fullName"`
	} `json:"contact"`
	Origin struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1" validate:"required" example:"Inglaterra 59"`
			AddressLine2 string `json:"addressLine2" example:"La Florida, Región Metropolitana, Chile"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceId" example:"BODEGA_2214"`
		} `json:"nodeInfo"`
	} `json:"origin"`
}

func (req CreateAccountRequest) Map() domain.Account {
	return domain.Account{
		Contact: domain.Contact{
			Email:      req.Contact.Email,
			NationalID: req.Contact.NationalID,
		},
		Origin: domain.Origin{
			AddressInfo: domain.AddressInfo{
				AddressLine1: req.Origin.AddressInfo.AddressLine1,
				AddressLine2: req.Origin.AddressInfo.AddressLine2,
			},
			NodeInfo: domain.NodeInfo{
				ReferenceID: domain.ReferenceID(req.Origin.NodeInfo.ReferenceID),
			},
		},
		Profiles: []domain.Profile{}, // Assuming profiles will be populated later
	}
}
