package response

import "transport-app/app/domain"

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
			County    string  `json:"county"`
			District  string  `json:"district"`
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
			Province  string  `json:"province"`
			State     string  `json:"state"`
			TimeZone  string  `json:"timeZone"`
			ZipCode   string  `json:"zipCode"`
		} `json:"addressInfo"`
		DeliveryInstructions string `json:"deliveryInstructions"`
		NodeInfo             struct {
			ReferenceID string `json:"referenceId"`
		} `json:"nodeInfo"`
	} `json:"destination"`
	Items []struct {
		Description string `json:"description"`
		Dimensions  struct {
			Depth  float64 `json:"depth"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
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
			County    string  `json:"county"`
			District  string  `json:"district"`
			Latitude  float32 `json:"latitude"`
			Longitude float32 `json:"longitude"`
			Province  string  `json:"province"`
			State     string  `json:"state"`
			TimeZone  string  `json:"timeZone"`
			ZipCode   string  `json:"zipCode"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceId"`
		} `json:"nodeInfo"`
	} `json:"origin"`
	Packages []struct {
		Dimensions struct {
			Depth  float64 `json:"depth"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
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
	Visits []struct {
		Date      string `json:"date"`
		TimeRange struct {
			EndTime   string `json:"endTime"`
			StartTime string `json:"startTime"`
		} `json:"timeRange"`
	} `json:"visits"`
}

func MapSearchOrdersResponse(orders []domain.Order) []SearchOrdersResponse {
	var responses []SearchOrdersResponse
	for _, order := range orders {
		response := SearchOrdersResponse{}
		response.
			withReferenceID(order.ReferenceID).
			withBusinessIdentifiers(order.BusinessIdentifiers).
			withOrderStatus(order.OrderStatus).
			withCollectAvailabilityDate(order.CollectAvailabilityDate).
			withOrigin(order.Origin).
			withDestination(order.Destination).
			withPromisedDate(order.PromisedDate).
			withVisits(order.Visits)
		responses = append(responses, response)
	}
	return responses
}

func (res *SearchOrdersResponse) withBusinessIdentifiers(businessIdentifiers domain.BusinessIdentifiers) *SearchOrdersResponse {
	res.BusinessIdentifiers.Commerce = businessIdentifiers.Commerce
	res.BusinessIdentifiers.Consumer = businessIdentifiers.Consumer
	return res
}

func (res *SearchOrdersResponse) withReferenceID(referenceID domain.ReferenceID) *SearchOrdersResponse {
	res.ReferenceID = string(referenceID)
	return res
}

func (res *SearchOrdersResponse) withOrderStatus(orderStatus domain.OrderStatus) *SearchOrdersResponse {
	res.OrderStatus.ID = orderStatus.ID
	res.OrderStatus.Status = orderStatus.Status
	res.OrderStatus.CreatedAt = orderStatus.CreatedAt
	return res
}

func (res *SearchOrdersResponse) withCollectAvailabilityDate(collectAvailabilityDate domain.CollectAvailabilityDate) *SearchOrdersResponse {
	res.CollectAvailabilityDate.Date = collectAvailabilityDate.Date
	res.CollectAvailabilityDate.TimeRange.EndTime = collectAvailabilityDate.TimeRange.EndTime
	res.CollectAvailabilityDate.TimeRange.StartTime = collectAvailabilityDate.TimeRange.StartTime
	return res
}

func (res *SearchOrdersResponse) withOrigin(origin domain.Origin) *SearchOrdersResponse {
	res.Origin.AddressInfo.AddressLine1 = origin.AddressInfo.AddressLine1
	res.Origin.AddressInfo.AddressLine2 = origin.AddressInfo.AddressLine2
	res.Origin.AddressInfo.AddressLine3 = origin.AddressInfo.AddressLine3
	res.Origin.AddressInfo.County = origin.AddressInfo.County
	res.Origin.AddressInfo.District = origin.AddressInfo.District
	res.Origin.AddressInfo.Province = origin.AddressInfo.Province
	res.Origin.AddressInfo.State = origin.AddressInfo.State
	res.Origin.AddressInfo.ZipCode = origin.AddressInfo.ZipCode
	res.Origin.AddressInfo.TimeZone = origin.AddressInfo.TimeZone
	res.Origin.AddressInfo.Latitude = origin.AddressInfo.Latitude
	res.Origin.AddressInfo.Longitude = origin.AddressInfo.Longitude
	res.Origin.NodeInfo.ReferenceID = string(origin.NodeInfo.ReferenceID)
	return res
}

func (res *SearchOrdersResponse) withDestination(destination domain.Destination) *SearchOrdersResponse {
	res.Destination.AddressInfo.AddressLine1 = destination.AddressInfo.AddressLine1
	res.Destination.AddressInfo.AddressLine2 = destination.AddressInfo.AddressLine2
	res.Destination.AddressInfo.AddressLine3 = destination.AddressInfo.AddressLine3
	res.Destination.AddressInfo.County = destination.AddressInfo.County
	res.Destination.AddressInfo.District = destination.AddressInfo.District
	res.Destination.AddressInfo.Province = destination.AddressInfo.Province
	res.Destination.AddressInfo.State = destination.AddressInfo.State
	res.Destination.AddressInfo.ZipCode = destination.AddressInfo.ZipCode
	res.Destination.AddressInfo.TimeZone = destination.AddressInfo.TimeZone
	res.Destination.AddressInfo.Latitude = destination.AddressInfo.Latitude
	res.Destination.AddressInfo.Longitude = destination.AddressInfo.Longitude
	res.Destination.DeliveryInstructions = destination.DeliveryInstructions
	res.Destination.NodeInfo.ReferenceID = string(destination.NodeInfo.ReferenceID)
	return res
}

func (res *SearchOrdersResponse) withPromisedDate(promisedDate domain.PromisedDate) *SearchOrdersResponse {
	res.PromisedDate.DateRange.StartDate = promisedDate.DateRange.StartDate
	res.PromisedDate.DateRange.EndDate = promisedDate.DateRange.EndDate
	res.PromisedDate.ServiceCategory = promisedDate.ServiceCategory
	res.PromisedDate.TimeRange.StartTime = promisedDate.TimeRange.StartTime
	res.PromisedDate.TimeRange.EndTime = promisedDate.TimeRange.EndTime
	return res
}

func (res *SearchOrdersResponse) withVisits(visits []domain.Visit) *SearchOrdersResponse {
	if res.Visits == nil {
		res.Visits = make([]struct {
			Date      string `json:"date"`
			TimeRange struct {
				EndTime   string `json:"endTime"`
				StartTime string `json:"startTime"`
			} `json:"timeRange"`
		}, 0)
	}

	res.Visits = res.Visits[:0]

	for _, visit := range visits {
		visitData := struct {
			Date      string `json:"date"`
			TimeRange struct {
				EndTime   string `json:"endTime"`
				StartTime string `json:"startTime"`
			} `json:"timeRange"`
		}{
			Date: visit.Date,
		}
		visitData.TimeRange.StartTime = visit.TimeRange.StartTime
		visitData.TimeRange.EndTime = visit.TimeRange.EndTime
		res.Visits = append(res.Visits, visitData)
	}

	return res
}
