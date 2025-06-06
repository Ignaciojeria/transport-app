package request

type RouteStartedRequest struct {
	StartedAt string `json:"startedAt" example:"2025-06-06T14:30:00Z"`
	Carrier   struct {
		Name       string `json:"name" example:"Transportes ABC"`
		NationalID string `json:"nationalID" example:"1234567890"`
	} `json:"carrier"`
	Driver struct {
		Email      string `json:"email" example:"juan@example.com"`
		NationalID string `json:"nationalID" example:"1234567890"`
	} `json:"driver"`
	Vehicle struct {
		Plate string `json:"plate" example:"ABC123"`
	} `json:"vehicle"`
	Route struct {
		Orders []struct {
			BusinessIdentifiers struct {
				Commerce string `json:"commerce"`
				Consumer string `json:"consumer"`
			} `json:"businessIdentifiers"`
			DeliveryUnits []struct {
				Items []struct {
					Sku string `json:"sku" example:"SKU123"`
				} `json:"items"`
				Lpn string `json:"lpn" example:"ABC123"`
			} `json:"deliveryUnits"`
			ReferenceID string `json:"referenceID"`
		} `json:"orders"`
		ReferenceID string `json:"referenceID"`
	} `json:"route"`
}
