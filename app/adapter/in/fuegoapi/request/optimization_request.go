package request

type OptimizationRequest struct {
	Vehicles []struct {
		Plate         string `json:"plate" example:"SERV-80" description:"Vehicle license plate or internal code"`
		StartLocation struct {
			Latitude  float64 `json:"latitude" example:"-33.45" description:"Starting point latitude"`
			Longitude float64 `json:"longitude" example:"-70.66" description:"Starting point longitude"`
		} `json:"startLocation"`
		EndLocation struct {
			Latitude  float64 `json:"latitude" example:"-33.45" description:"Ending point latitude"`
			Longitude float64 `json:"longitude" example:"-70.66" description:"Ending point longitude"`
		} `json:"endLocation"`
		Skills     []string `json:"skills" description:"Vehicle capabilities such as size or equipment requirements. eg: XL, heavy, etc"`
		TimeWindow struct {
			Start string `json:"start" example:"08:00" description:"Time window start (24h format)"`
			End   string `json:"end" example:"18:00" description:"Time window end (24h format)"`
		} `json:"timeWindow"`
		Capacity struct {
			Insurance             int `json:"insurance" example:"100000" description:"Maximum insurance value the vehicle can carry (CLP,MXN,PEN)"`
			Weight                int `json:"weight" example:"1000" description:"Maximum weight in kilograms"`
			DeliveryUnitsQuantity int `json:"deliveryUnitsQuantity" example:"50" description:"Maximum number of delivery units the vehicle can carry"`
		} `json:"capacity"`
	} `json:"vehicles"`
	Visits []struct {
		PickupLocation struct {
			Latitude  float64 `json:"latitude" example:"-33.45" description:"Pickup point latitude"`
			Longitude float64 `json:"longitude" example:"-70.66" description:"Pickup point longitude"`
		} `json:"pickupLocation"`
		DispatchLocation struct {
			Latitude  float64 `json:"latitude" example:"-33.45" description:"Dispatch point latitude"`
			Longitude float64 `json:"longitude" example:"-70.66" description:"Dispatch point longitude"`
		} `json:"dispatchLocation"`
		CapacityUsage struct {
			Insurance             int `json:"insurance" example:"50000" description:"Insurance value of the delivery units"`
			Weight                int `json:"weight" example:"500" description:"Total weight of the delivery units in kilograms"`
			DeliveryUnitsQuantity int `json:"deliveryUnitsQuantity" example:"25" description:"Number of delivery units in this visit"`
		} `json:"capacityUsage"`
		Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
		TimeWindow struct {
			Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
			End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
		} `json:"timeWindow"`
		ServiceTime int `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
		Orders      []struct {
			DeliveryUnits []struct {
				Items []struct {
					Sku string `json:"sku" example:"SKU123" description:"Stock keeping unit identifier"`
				} `json:"items"`
				Lpn string `json:"lpn" example:"LPN456" description:"License plate number of the delivery unit"`
			} `json:"deliveryUnits"`
			ReferenceID string `json:"referenceID" example:"ORD789" description:"Unique identifier for the order"`
		} `json:"orders"`
	} `json:"visits"`
}
