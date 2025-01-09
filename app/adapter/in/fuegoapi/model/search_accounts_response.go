package model

type SearchAccountsResponse struct {
	Contact struct {
		Email      string `json:"email"`
		FullName   string `json:"fullName"`
		NationalID string `json:"nationalID"`
		Phone      string `json:"phone"`
	} `json:"contact"`
	Origin struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceId"`
		} `json:"nodeInfo"`
	} `json:"origin"`
}
