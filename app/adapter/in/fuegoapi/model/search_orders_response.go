package model

type SearchOrdersResponse struct {
	ReferenceID         string `json:"referenceID"`
	BusinessIdentifiers struct {
		Commerce string `json:"commerce"`
		Consumer string `json:"consumer"`
	} `json:"businessIdentifiers"`
	OrderStatus struct {
		ID        int64
		Status    string `json:"status"`
		CreatedAt string `json:"createdAt"`
	} `json:"orderStatus"`
	CollectAvailabilityDate struct {
		Date      string `json:"date"`
		TimeRange struct {
			EndTime   string `json:"endTime"`
			StartTime string `json:"startTime"`
		} `json:"timeRange"`
	} `json:"collectAvailabilityDate"`
	Destination struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
			AddressLine3 string `json:"addressLine3"`
			Contact      struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				Documents  []struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"documents"`
				FullName string `json:"fullName"`
			} `json:"contact"`
			County    string `json:"county"`
			District  string `json:"district"`
			Latitude  int    `json:"latitude"`
			Longitude int    `json:"longitude"`
			Province  string `json:"province"`
			State     string `json:"state"`
			TimeZone  string `json:"timeZone"`
			ZipCode   string `json:"zipCode"`
		} `json:"addressInfo"`
		DeliveryInstructions string `json:"deliveryInstructions"`
		NodeInfo             struct {
			ReferenceID string `json:"referenceId"`
		} `json:"nodeInfo"`
	} `json:"destination"`
	Items []struct {
		Description string `json:"description"`
		Dimensions  struct {
			Depth  int    `json:"depth"`
			Height int    `json:"height"`
			Unit   string `json:"unit"`
			Width  int    `json:"width"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string `json:"currency"`
			UnitValue int    `json:"unitValue"`
		} `json:"insurance"`
		LogisticCondition string `json:"logisticCondition"`
		Quantity          struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		ReferenceID string `json:"referenceId"`
		Weight      struct {
			Unit  string `json:"unit"`
			Value int    `json:"value"`
		} `json:"weight"`
	} `json:"items"`
	OrderType struct {
		Description string `json:"description"`
		Type        string `json:"type"`
	} `json:"orderType"`
	Origin struct {
		AddressInfo struct {
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
			AddressLine3 string `json:"addressLine3"`
			Contact      struct {
				Email      string `json:"email"`
				Phone      string `json:"phone"`
				NationalID string `json:"nationalID"`
				Documents  []struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"documents"`
				FullName string `json:"fullName"`
			} `json:"contact"`
			County    string `json:"county"`
			District  string `json:"district"`
			Latitude  int    `json:"latitude"`
			Longitude int    `json:"longitude"`
			Province  string `json:"province"`
			State     string `json:"state"`
			TimeZone  string `json:"timeZone"`
			ZipCode   string `json:"zipCode"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceId"`
		} `json:"nodeInfo"`
	} `json:"origin"`
	Packages []struct {
		Dimensions struct {
			Depth  int    `json:"depth"`
			Height int    `json:"height"`
			Unit   string `json:"unit"`
			Width  int    `json:"width"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string `json:"currency"`
			UnitValue int    `json:"unitValue"`
		} `json:"insurance"`
		ItemReferences []struct {
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			ReferenceID string `json:"referenceId"`
		} `json:"itemReferences"`
		Lpn         string `json:"lpn"`
		PackageType string `json:"packageType"`
		Weight      struct {
			Unit  string `json:"unit"`
			Value int    `json:"value"`
		} `json:"weight"`
	} `json:"packages"`
	PromisedDate struct {
		DateRange struct {
			EndDate   string `json:"endDate"`
			StartDate string `json:"startDate"`
		} `json:"dateRange"`
		ServiceCategory string `json:"serviceCategory"`
		TimeRange       struct {
			EndTime   string `json:"endTime"`
			StartTime string `json:"startTime"`
		} `json:"timeRange"`
	} `json:"promisedDate"`
	References []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"references"`
	TransportRequirements []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"transportRequirements"`
	Visit struct {
		Date      string `json:"date"`
		TimeRange struct {
			EndTime   string `json:"endTime"`
			StartTime string `json:"startTime"`
		} `json:"timeRange"`
	} `json:"visit"`
}
