package request

type ConfirmDeliveriesRequest struct {
	Carrier struct {
		Name       string `json:"name"`
		NationalID string `json:"nationalID"`
	} `json:"carrier"`
	Driver struct {
		Email      string `json:"email"`
		NationalID string `json:"nationalID"`
	} `json:"driver"`
	Routes []struct {
		Orders []struct {
			BusinessIdentifiers struct {
				Commerce string `json:"commerce"`
				Consumer string `json:"consumer"`
			} `json:"businessIdentifiers"`
			Delivery struct {
				Failure struct {
					Detail      string `json:"detail"`
					Reason      string `json:"reason"`
					ReferenceID string `json:"referenceID"`
				} `json:"failure"`
				HandledAt string `json:"handledAt"`
				Location  struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"location"`
			} `json:"delivery"`
			EvidencePhotos []struct {
				TakenAt string `json:"takenAt"`
				Type    string `json:"type"`
				URL     string `json:"url"`
			} `json:"evidencePhotos"`
			DeliveryUnits []struct {
				Items []struct {
					Sku string `json:"sku"`
				} `json:"items"`
				Lpn string `json:"lpn"`
			} `json:"deliveryUnits"`
			Recipient struct {
				FullName   string `json:"fullName"`
				NationalID string `json:"nationalID"`
			} `json:"recipient"`
			ReferenceID string `json:"referenceID"`
		} `json:"orders"`
		ReferenceID string `json:"referenceID"`
	} `json:"routes"`
	Vehicle struct {
		Plate string `json:"plate"`
	} `json:"vehicle"`
}
