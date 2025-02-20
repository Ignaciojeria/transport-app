package response

import "transport-app/app/domain"

type SearchAccountResponse struct {
	ReferenceID string `json:"referenceID"`
	Contact     struct {
		Email      string `json:"email"`
		FullName   string `json:"fullName"`
		NationalID string `json:"nationalID"`
		Phone      string `json:"phone"`
	} `json:"contact"`
	OriginNode struct {
		ReferenceID string `json:"referenceID"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		AddressInfo struct {
			RawAddress string  `json:"rawAddress"`
			Latitude   float64 `json:"latitude"`
			Longitude  float64 `json:"longitude"`
		} `json:"addressInfo"`
	} `json:"originNode"`
}

func MapSearchAccountOperatorResponse(operator domain.Operator) SearchAccountResponse {
	return SearchAccountResponse{
		ReferenceID: operator.ReferenceID,
		Contact: struct {
			Email      string "json:\"email\""
			FullName   string "json:\"fullName\""
			NationalID string "json:\"nationalID\""
			Phone      string "json:\"phone\""
		}{
			Email:      operator.Contact.Email,
			FullName:   operator.Contact.FullName,
			NationalID: operator.Contact.NationalID,
			Phone:      operator.Contact.Phone,
		},
		OriginNode: struct {
			ReferenceID string "json:\"referenceID\""
			Name        string "json:\"name\""
			Type        string "json:\"type\""
			AddressInfo struct {
				RawAddress string  "json:\"rawAddress\""
				Latitude   float64 "json:\"latitude\""
				Longitude  float64 "json:\"longitude\""
			} "json:\"addressInfo\""
		}{
			ReferenceID: string(operator.OriginNode.ReferenceID),
			Name:        operator.OriginNode.Name,
			Type:        operator.Type,
			AddressInfo: struct {
				RawAddress string  "json:\"rawAddress\""
				Latitude   float64 "json:\"latitude\""
				Longitude  float64 "json:\"longitude\""
			}{
				RawAddress: operator.OriginNode.AddressInfo.RawAddress(),
				Longitude:  operator.OriginNode.AddressInfo.Location[0],
				Latitude:   operator.OriginNode.AddressInfo.Location[1],
			},
		},
	}
}
