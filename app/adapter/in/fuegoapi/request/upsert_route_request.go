package request

type UpsertRouteRequest struct {
	ReferenceID string `json:"referenceID" example:"ROUTE-001"`
	Plan        struct {
		ReferenceID string `json:"referenceID" example:"PLAN-001"`
	} `json:"plan"`
	Vehicle struct {
		Plate string `json:"plate" example:"ABCD12"`
	} `json:"vehicle"`
	Geometry struct {
		Encoding string `json:"encoding" example:"polyline"`
		Type     string `json:"type" example:"linestring"`
		Value    string `json:"value" example:"_p~iF~ps|U_ulLnnqC_mqNvxq@"`
	} `json:"geometry"`
	Visits []struct {
		Type        string `json:"type" example:"delivery"`
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1" example:"Av. Providencia 1234"`
			AddressLine2 string `json:"addressLine2" example:"Oficina 567"`
			Contact      struct {
				Email      string `json:"email" example:"cliente@ejemplo.com"`
				FullName   string `json:"fullName" example:"Juan Pérez"`
				NationalID string `json:"nationalID" example:"12345678-9"`
				Phone      string `json:"phone" example:"+56912345678"`
			} `json:"contact"`
			Coordinates struct {
				Latitude  float64 `json:"latitude" example:"-33.5147889"`
				Longitude float64 `json:"longitude" example:"-70.6130425"`
			} `json:"coordinates"`
			PoliticalArea struct {
				Code     string `json:"code" example:"cl-rm-providencia"`
				District string `json:"district" example:"providencia"`
				Province string `json:"province" example:"santiago"`
				State    string `json:"state" example:"region metropolitana de santiago"`
			} `json:"politicalArea"`
			ZipCode string `json:"zipCode" example:"7500000"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID" example:"NODE-001"`
		} `json:"nodeInfo"`
		DeliveryInstructions string `json:"deliveryInstructions" example:"Entregar en recepción"`
		SequenceNumber       int    `json:"sequenceNumber" example:"1"`
		Orders               []struct {
			DeliveryUnits []struct {
				Items []struct {
					Sku string `json:"sku" example:"SKU-123456"`
				} `json:"items"`
				Lpn string `json:"lpn" example:"LPN-789012"`
			} `json:"deliveryUnits"`
			ReferenceID string `json:"referenceID" example:"ORDER-001"`
		} `json:"orders"`
	} `json:"visits"`
	CreatedAt string `json:"createdAt" example:"2025-01-15T10:30:00Z"`
}
