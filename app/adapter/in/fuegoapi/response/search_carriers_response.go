package response

import "transport-app/app/domain"

type SearchCarriersResponse struct {
	ReferenceID string `json:"referenceID"`
	Name        string `json:"name"`
	NationalID  string `json:"nationalID"`
}

func MapSearchCarriersResponse(carriers []domain.Carrier) []SearchCarriersResponse {
	response := make([]SearchCarriersResponse, len(carriers))
	for i, carrier := range carriers {
		response[i] = SearchCarriersResponse{
			Name:       carrier.Name,
			NationalID: carrier.NationalID,
		}
	}
	return response
}
