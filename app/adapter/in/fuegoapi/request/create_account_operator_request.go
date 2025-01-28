package request

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
