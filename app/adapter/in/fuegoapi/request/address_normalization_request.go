package request

import (
	"transport-app/app/domain"

	"github.com/paulmach/orb"
)

type AddressNormalizationRequest struct {
	ProviderInput struct {
		DisplayName string  `json:"displayName"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
	} `json:"providerInput"`
	UserInput struct {
		AddressLine1 string `json:"addressLine1"`
		AddressLine2 string `json:"addressLine2"`
		State        string `json:"state"`
		District     string `json:"district"`
	} `json:"userInput"`
}

func (req AddressNormalizationRequest) Map() (domain.AddressInfo, domain.AddressInfo) {
	return domain.AddressInfo{
			AddressLine1: req.UserInput.AddressLine1,
			//AddressLine2: req.UserInput.AddressLine2,
			State:        req.UserInput.State,
			District:     req.UserInput.District,
		}, domain.AddressInfo{
			AddressLine1: req.ProviderInput.DisplayName,
			Location: orb.Point{
				req.ProviderInput.Longitude,
				req.ProviderInput.Latitude,
			},
		}
}
