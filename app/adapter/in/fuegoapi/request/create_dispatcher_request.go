package request

import "transport-app/app/domain"

type CreateDispatcherRequest struct {
	Contact struct {
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		NationalID string `json:"nationalID"`
		FullName   string `json:"fullName"`
	} `json:"contact"`
	NodeReferenceIDs []string `json:"nodeReferenceIDs"`
}

func (req CreateDispatcherRequest) Map() domain.Account {
	return domain.Account{
		Contact: domain.Contact{
			Email:      req.Contact.Email,
			NationalID: req.Contact.NationalID,
		},
		Profiles: []domain.Profile{}, // Assuming profiles will be populated later
	}
}
