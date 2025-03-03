package response

import "transport-app/app/domain"

// AddressNormalizationResponse representa la dirección normalizada devuelta en la respuesta.
type AddressNormalizationResponse struct {
	ProviderAddress string  `json:"providerAddress"`
	AddressLine1    string  `json:"addressLine1"`
	AddressLine2    string  `json:"addressLine2,omitempty"`
	District        string  `json:"district"`
	Province        string  `json:"province"`
	State           string  `json:"state"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}

// MapAddressNormalizationResponse convierte una estructura domain.AddressInfo en AddressNormalizationResponse.
func MapAddressNormalizationResponse(addressInfo domain.AddressInfo) AddressNormalizationResponse {
	return AddressNormalizationResponse{
		ProviderAddress: addressInfo.ProviderAddress,
		AddressLine1:    addressInfo.AddressLine1,
		AddressLine2:    addressInfo.AddressLine2,
		District:        addressInfo.District,
		Province:        addressInfo.Province,
		State:           addressInfo.State,
		Latitude:        addressInfo.Location.Lat(),
		Longitude:       addressInfo.Location.Lon(),
	}
}
