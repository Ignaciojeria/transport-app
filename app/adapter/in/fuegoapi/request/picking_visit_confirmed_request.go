package request

type PickingVisitConfirmedRequest struct {
	Carrier struct {
		Name       string `json:"name"`
		NationalID string `json:"nationalID"`
	} `json:"carrier"`
	Vehicle struct {
		Plate string `json:"plate"`
	} `json:"vehicle"`
	Driver struct {
		Email      string `json:"email"`
		NationalID string `json:"nationalID"`
	} `json:"driver"`
	Visit struct {
		Comment  string `json:"comment"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID"`
			SellerID    string `json:"sellerID"`
		} `json:"nodeInfo"`
		Contact struct {
			Email      string `json:"email"`
			FullName   string `json:"fullName"`
			NationalID string `json:"nationalID"`
			Phone      string `json:"phone"`
		} `json:"contact"`
		Coordinates struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"coordinates"`
		VisitDate        string `json:"visitDate" example:"2025-06-19T10:00:00Z" description:"Visit date"`
		ConfirmationDate string `json:"confirmationDate" example:"2025-06-19T10:00:00Z" description:"Confirmation date"`
	} `json:"visit"`
	Orders []struct {
		DeliveryUnits []struct {
			Price int `json:"price"`
			Items []struct {
				Sku string `json:"sku"`
			} `json:"items"`
			Lpn    string `json:"lpn"`
			Volume int    `json:"volume"`
			Weight int    `json:"weight"`
		} `json:"deliveryUnits"`
		ReferenceID string `json:"referenceID"`
	} `json:"orders"`
}
