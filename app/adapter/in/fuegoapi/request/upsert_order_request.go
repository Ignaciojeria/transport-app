package request

import (
	"time"
	"transport-app/app/domain"

	"github.com/paulmach/orb"
)

type UpsertOrderRequest struct {
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
			County    string  `json:"county"`
			District  string  `json:"district"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Province  string  `json:"province"`
			State     string  `json:"state"`
			TimeZone  string  `json:"timeZone"`
			ZipCode   string  `json:"zipCode"`
		} `json:"addressInfo"`
		DeliveryInstructions string `json:"deliveryInstructions"`
		NodeInfo             struct {
			ReferenceID string `json:"referenceId"`
			Name        string `json:"name"`
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
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		} `json:"insurance"`
		LogisticCondition string `json:"logisticCondition"`
		Quantity          struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		ReferenceID string `json:"referenceId"`
		Weight      struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
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
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Province  string  `json:"province"`
			State     string  `json:"state"`
			TimeZone  string  `json:"timeZone"`
			ZipCode   string  `json:"zipCode"`
		} `json:"addressInfo"`
		NodeInfo struct {
			ReferenceID string `json:"referenceId"`
			Name        string `json:"name"`
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
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		} `json:"insurance"`
		ItemReferences []struct {
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			ReferenceID string `json:"referenceId"`
		} `json:"itemReferences"`
		Lpn    string `json:"lpn"`
		Weight struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
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
}

func (req UpsertOrderRequest) Map() domain.Order {
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
		DeliveryInstructions:    req.Destination.DeliveryInstructions,
		//Visits:                  req.mapVisit(),
		TransportRequirements: req.mapReferences(req.TransportRequirements),
	}
}

func (req UpsertOrderRequest) mapOrderType() domain.OrderType {
	return domain.OrderType{
		Type:        req.OrderType.Type,
		Description: req.OrderType.Description,
	}
}

func (req UpsertOrderRequest) mapReferences(refs []struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}) []domain.Reference {
	mapped := make([]domain.Reference, len(refs))
	for i, ref := range refs {
		mapped[i] = domain.Reference{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return mapped
}

func (req UpsertOrderRequest) mapOrigin() domain.NodeInfo {
	nodeInfo := req.mapNodeInfo(req.Origin.NodeInfo)
	nodeInfo.AddressInfo = req.mapAddressInfo(req.Origin.AddressInfo)
	return nodeInfo
}

func (req UpsertOrderRequest) mapDestination() domain.NodeInfo {
	nodeInfo := req.mapNodeInfo(req.Destination.NodeInfo)
	nodeInfo.AddressInfo = req.mapAddressInfo(req.Destination.AddressInfo)
	return nodeInfo
}

func (req UpsertOrderRequest) mapNodeInfo(nodeInfo struct {
	ReferenceID string `json:"referenceId"`
	Name        string `json:"name"`
}) domain.NodeInfo {
	return domain.NodeInfo{
		ReferenceID: domain.ReferenceID(nodeInfo.ReferenceID),
		Name:        nodeInfo.Name,
	}
}

func (req UpsertOrderRequest) mapAddressInfo(addressInfo struct {
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
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Province  string  `json:"province"`
	State     string  `json:"state"`
	TimeZone  string  `json:"timeZone"`
	ZipCode   string  `json:"zipCode"`
}) domain.AddressInfo {
	return domain.AddressInfo{
		Contact: domain.Contact{
			FullName:   addressInfo.Contact.FullName,
			Email:      addressInfo.Contact.Email,
			Phone:      addressInfo.Contact.Phone,
			NationalID: addressInfo.Contact.NationalID,
			Documents:  req.mapDocuments(addressInfo.Contact.Documents),
		},
		State:        addressInfo.State,
		County:       addressInfo.County,
		Province:     addressInfo.Province,
		District:     addressInfo.District,
		AddressLine1: addressInfo.AddressLine1,
		AddressLine2: addressInfo.AddressLine2,
		AddressLine3: addressInfo.AddressLine3,
		Location: orb.Point{
			addressInfo.Longitude, // orb.Point espera [lon, lat]
			addressInfo.Latitude,
		},
		ZipCode:  addressInfo.ZipCode,
		TimeZone: addressInfo.TimeZone,
	}
}

func (req UpsertOrderRequest) mapDocuments(docs []struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}) []domain.Document {
	mapped := make([]domain.Document, len(docs))
	for i, doc := range docs {
		mapped[i] = domain.Document{
			Type:  doc.Type,
			Value: doc.Value,
		}
	}
	return mapped
}

func (req UpsertOrderRequest) mapItems() []domain.Item {
	mapped := make([]domain.Item, len(req.Items))
	for i, item := range req.Items {
		mapped[i] = domain.Item{
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
				Height: item.Dimensions.Height,
				Width:  item.Dimensions.Width,
				Depth:  item.Dimensions.Depth,
				Unit:   item.Dimensions.Unit,
			},
			Weight: domain.Weight{
				Unit:  item.Weight.Unit,
				Value: item.Weight.Value,
			},
		}
	}
	return mapped
}

func (req UpsertOrderRequest) mapPackages() []domain.Package {
	mapped := make([]domain.Package, len(req.Packages))
	for i, pkg := range req.Packages {
		mapped[i] = domain.Package{
			Lpn: pkg.Lpn,
			Dimensions: domain.Dimensions{
				Height: pkg.Dimensions.Height,
				Width:  pkg.Dimensions.Width,
				Depth:  pkg.Dimensions.Depth,
				Unit:   pkg.Dimensions.Unit,
			},
			Weight: domain.Weight{
				Unit:  pkg.Weight.Unit,
				Value: pkg.Weight.Value,
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

func (req UpsertOrderRequest) mapItemReferences(itemReferences []struct {
	Quantity struct {
		QuantityNumber int    `json:"quantityNumber"`
		QuantityUnit   string `json:"quantityUnit"`
	} `json:"quantity"`
	ReferenceID string `json:"referenceId"`
}) []domain.ItemReference {
	mapped := make([]domain.ItemReference, len(itemReferences))
	for i, ref := range itemReferences {
		mapped[i] = domain.ItemReference{
			ReferenceID: domain.ReferenceID(ref.ReferenceID),
			Quantity: domain.Quantity{
				QuantityNumber: ref.Quantity.QuantityNumber,
				QuantityUnit:   ref.Quantity.QuantityUnit,
			},
		}
	}
	return mapped
}
func (req UpsertOrderRequest) mapCollectAvailabilityDate() domain.CollectAvailabilityDate {
	date, err := time.Parse("2006-01-02", req.CollectAvailabilityDate.Date)
	if err != nil {
		// Dependiendo de tu manejo de errores, podrías:
		// 1. Retornar un zero value time.Time
		// 2. Usar un valor por defecto
		// 3. Propagar el error (requeriría cambiar la firma del método)
		date = time.Time{} // zero value como ejemplo
	}

	return domain.CollectAvailabilityDate{
		Date: date,
		TimeRange: domain.TimeRange{
			StartTime: req.CollectAvailabilityDate.TimeRange.StartTime,
			EndTime:   req.CollectAvailabilityDate.TimeRange.EndTime,
		},
	}
}

func (req UpsertOrderRequest) mapPromisedDate() domain.PromisedDate {
	startDate, err := time.Parse("2006-01-02", req.PromisedDate.DateRange.StartDate)
	if err != nil {
		startDate = time.Time{}
	}

	endDate, err := time.Parse("2006-01-02", req.PromisedDate.DateRange.EndDate)
	if err != nil {
		endDate = time.Time{}
	}

	return domain.PromisedDate{
		DateRange: domain.DateRange{
			StartDate: startDate,
			EndDate:   endDate,
		},
		TimeRange: domain.TimeRange{
			StartTime: req.PromisedDate.TimeRange.StartTime,
			EndTime:   req.PromisedDate.TimeRange.EndTime,
		},
		ServiceCategory: req.PromisedDate.ServiceCategory,
	}
}
