package model

import "transport-app/app/domain"

type CreateOrderRequest struct {
	ReferenceID             string `json:"referenceID" validate:"required"`
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

func (req CreateOrderRequest) Map() domain.Order {
	return domain.Order{
		ReferenceID:             domain.ReferenceID(req.ReferenceID),
		OrderType:               req.mapOrderType(),
		References:              req.mapReferences(req.References),
		Origin:                  req.mapOrigin(),
		Destination:             req.mapDestination(),
		Items:                   req.mapItems(),
		Packages:                req.mapPackages(),
		CollectAvailabilityDate: req.mapCollectAvailabilityDate(),
		PromisedDate:            req.mapPromisedDate(),
		Visit:                   req.mapVisit(),
		TransportRequirements:   req.mapReferences(req.TransportRequirements),
	}
}

func (req CreateOrderRequest) mapOrderType() domain.OrderType {
	return domain.OrderType{
		Type:        req.OrderType.Type,
		Description: req.OrderType.Description,
	}
}

func (req CreateOrderRequest) mapReferences(refs []struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}) []domain.References {
	mapped := make([]domain.References, len(refs))
	for i, ref := range refs {
		mapped[i] = domain.References{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return mapped
}

func (req CreateOrderRequest) mapOrigin() domain.Origin {
	return domain.Origin{
		NodeInfo:    req.mapNodeInfo(req.Origin.NodeInfo),
		AddressInfo: req.mapAddressInfo(req.Origin.AddressInfo),
	}
}

func (req CreateOrderRequest) mapDestination() domain.Destination {
	return domain.Destination{
		DeliveryInstructions: req.Destination.DeliveryInstructions,
		NodeInfo:             req.mapNodeInfo(req.Destination.NodeInfo),
		AddressInfo:          req.mapAddressInfo(req.Destination.AddressInfo),
	}
}

func (req CreateOrderRequest) mapNodeInfo(nodeInfo struct {
	ReferenceID string `json:"referenceId"`
}) domain.NodeInfo {
	return domain.NodeInfo{
		ReferenceID: domain.ReferenceID(nodeInfo.ReferenceID),
	}
}

func (req CreateOrderRequest) mapAddressInfo(addressInfo struct {
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
}) domain.AddressInfo {
	return domain.AddressInfo{
		Contact: domain.Contact{
			FullName:  addressInfo.Contact.FullName,
			Email:     addressInfo.Contact.Email,
			Phone:     addressInfo.Contact.Phone,
			Documents: req.mapDocuments(addressInfo.Contact.Documents),
		},
		State:        addressInfo.State,
		County:       addressInfo.County,
		Province:     addressInfo.Province,
		District:     addressInfo.District,
		AddressLine1: addressInfo.AddressLine1,
		AddressLine2: addressInfo.AddressLine2,
		AddressLine3: addressInfo.AddressLine3,
		Latitude:     float64(addressInfo.Latitude),
		Longitude:    float64(addressInfo.Longitude),
		ZipCode:      addressInfo.ZipCode,
		TimeZone:     addressInfo.TimeZone,
	}
}

func (req CreateOrderRequest) mapDocuments(docs []struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}) []domain.Documents {
	mapped := make([]domain.Documents, len(docs))
	for i, doc := range docs {
		mapped[i] = domain.Documents{
			Type:  doc.Type,
			Value: doc.Value,
		}
	}
	return mapped
}

func (req CreateOrderRequest) mapItems() []domain.Items {
	mapped := make([]domain.Items, len(req.Items))
	for i, item := range req.Items {
		mapped[i] = domain.Items{
			ReferenceID:       domain.ReferenceID(item.ReferenceID),
			LogisticCondition: item.LogisticCondition,
			Quantity: domain.Quantity{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
			Insurance: domain.Insurance{
				Currency:  item.Insurance.Currency,
				UnitValue: item.Insurance.UnitValue,
			},
			Description: item.Description,
			Dimensions: domain.Dimensions{
				Height: float64(item.Dimensions.Height),
				Width:  float64(item.Dimensions.Width),
				Depth:  float64(item.Dimensions.Depth),
				Unit:   item.Dimensions.Unit,
			},
			Weight: domain.Weight{
				Unit:  item.Weight.Unit,
				Value: float64(item.Weight.Value),
			},
		}
	}
	return mapped
}

func (req CreateOrderRequest) mapPackages() []domain.Packages {
	mapped := make([]domain.Packages, len(req.Packages))
	for i, pkg := range req.Packages {
		mapped[i] = domain.Packages{
			Lpn:         pkg.Lpn,
			PackageType: pkg.PackageType,
			Dimensions: domain.Dimensions{
				Height: float64(pkg.Dimensions.Height),
				Width:  float64(pkg.Dimensions.Width),
				Depth:  float64(pkg.Dimensions.Depth),
				Unit:   pkg.Dimensions.Unit,
			},
			Weight: domain.Weight{
				Unit:  pkg.Weight.Unit,
				Value: float64(pkg.Weight.Value),
			},
			Insurance: domain.Insurance{
				Currency:  pkg.Insurance.Currency,
				UnitValue: pkg.Insurance.UnitValue,
			},
			ItemReferences: req.mapItemReferences(pkg.ItemReferences),
		}
	}
	return mapped
}

func (req CreateOrderRequest) mapItemReferences(itemReferences []struct {
	Quantity struct {
		QuantityNumber int    `json:"quantityNumber"`
		QuantityUnit   string `json:"quantityUnit"`
	} `json:"quantity"`
	ReferenceID string `json:"referenceId"`
}) []domain.ItemReferences {
	mapped := make([]domain.ItemReferences, len(itemReferences))
	for i, ref := range itemReferences {
		mapped[i] = domain.ItemReferences{
			ReferenceID: domain.ReferenceID(ref.ReferenceID),
			Quantity: domain.Quantity{
				QuantityNumber: ref.Quantity.QuantityNumber,
				QuantityUnit:   ref.Quantity.QuantityUnit,
			},
		}
	}
	return mapped
}

func (req CreateOrderRequest) mapCollectAvailabilityDate() domain.CollectAvailabilityDate {
	return domain.CollectAvailabilityDate{
		Date: req.CollectAvailabilityDate.Date,
		TimeRange: domain.TimeRange{
			StartTime: req.CollectAvailabilityDate.TimeRange.StartTime,
			EndTime:   req.CollectAvailabilityDate.TimeRange.EndTime,
		},
	}
}

func (req CreateOrderRequest) mapPromisedDate() domain.PromisedDate {
	return domain.PromisedDate{
		DateRange: domain.DateRange{
			StartDate: req.PromisedDate.DateRange.StartDate,
			EndDate:   req.PromisedDate.DateRange.EndDate,
		},
		TimeRange: domain.TimeRange{
			StartTime: req.PromisedDate.TimeRange.StartTime,
			EndTime:   req.PromisedDate.TimeRange.EndTime,
		},
		ServiceCategory: req.PromisedDate.ServiceCategory,
	}
}

func (req CreateOrderRequest) mapVisit() domain.Visit {
	return domain.Visit{
		Date: req.Visit.Date,
		TimeRange: domain.TimeRange{
			StartTime: req.Visit.TimeRange.StartTime,
			EndTime:   req.Visit.TimeRange.EndTime,
		},
	}
}
