package request

import "transport-app/app/domain"

type CreateAccountOperatorRequest struct {
	ReferenceID           string `json:"referenceID"`
	OriginNodeReferenceID string `json:"originNodeReferenceID"`
	Contact               struct {
		Email      string `json:"email"`
		FullName   string `json:"fullName"`
		NationalID string `json:"nationalID"`
		Phone      string `json:"phone"`
	} `json:"contact"`
}

func (r CreateAccountOperatorRequest) Map() domain.Operator {
	return domain.Operator{
		//ReferenceID: r.ReferenceID,
		OriginNode: domain.NodeInfo{
			ReferenceID: domain.ReferenceID(r.OriginNodeReferenceID),
		},
		Contact: domain.Contact{
			PrimaryEmail: r.Contact.Email,
			FullName:     r.Contact.FullName,
			NationalID:   r.Contact.NationalID,
			PrimaryPhone: r.Contact.Phone,
		},
	}
}
