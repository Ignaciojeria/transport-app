package request

type OptimizePickingAndDeliveryRequest struct {
	NodeInfo struct {
		ReferenceID string `json:"referenceID"`
	} `json:"nodeInfo"`
	Container struct {
		Lpn string `json:"lpn" example:"LPN456" description:"License plate number of the container"`
	} `json:"container"`
	Carrier struct {
		Name       string `json:"name"`
		NationalID string `json:"nationalID"`
	} `json:"carrier"`
	Driver struct {
		Email      string `json:"email"`
		NationalID string `json:"nationalID"`
	} `json:"driver"`
	StartLocation struct {
		Latitude  float64 `json:"latitude" example:"-33.45" description:"Starting point latitude"`
		Longitude float64 `json:"longitude" example:"-70.66" description:"Starting point longitude"`
		NodeInfo  struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"startLocation"`
	EndLocation struct {
		Latitude  float64 `json:"latitude" example:"-33.45" description:"Ending point latitude"`
		Longitude float64 `json:"longitude" example:"-70.66" description:"Ending point longitude"`
		NodeInfo  struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"endLocation"`
	Visits []struct {
		Delivery struct {
			Coordinates struct {
				Latitude  float64 `json:"latitude" example:"-33.45" description:"Pickup point latitude"`
				Longitude float64 `json:"longitude" example:"-70.66" description:"Pickup point longitude"`
			} `json:"coordinates"`
			ServiceTime int64 `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
			Contact     struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				FullName   string `json:"fullName"`
			} `json:"contact"`
		} `json:"delivery"`
		Orders []struct {
			DeliveryUnits []struct {
				Items []struct {
					Sku string `json:"sku" example:"SKU123" description:"Stock keeping unit identifier"`
				} `json:"items"`
				Insurance int64  `json:"insurance" example:"10000" description:"Insurance value of the delivery unit"`
				Volume    int64  `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
				Weight    int64  `json:"weight" example:"1000" description:"Weight of the delivery unit in grams"`
				Lpn       string `json:"lpn" example:"LPN456" description:"License plate number of the delivery unit"`
			} `json:"deliveryUnits"`
			ReferenceID string `json:"referenceID" example:"ORD789" description:"Unique identifier for the order"`
		} `json:"orders"`
	} `json:"visits"`
}
