package request

import "transport-app/app/domain/optimization"

type OptimizeFleetRequest struct {
	PlanReferenceID string `json:"planReferenceID"`
	Vehicles        []struct {
		Plate         string `json:"plate" example:"SERV-80" description:"Vehicle license plate or internal code"`
		StartLocation struct {
			AddressInfo struct {
				AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
				AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
				Contact      struct {
					Email      string `json:"email"`
					Phone      string `json:"phone"`
					NationalID string `json:"nationalID"`
					FullName   string `json:"fullName"`
				} `json:"contact"`
				Coordinates struct {
					Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Starting point latitude"`
					Longitude float64 `json:"longitude" example:"-70.6130425" description:"Starting point longitude"`
				} `json:"coordinates"`
				PoliticalArea struct {
					Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
					District string `json:"district" example:"la florida" description:"District name"`
					Province string `json:"province" example:"santiago" description:"Province name"`
					State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
				} `json:"politicalArea"`
				ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
			} `json:"addressInfo"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
		} `json:"startLocation"`
		EndLocation struct {
			AddressInfo struct {
				AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
				AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
				Contact      struct {
					Email      string `json:"email"`
					Phone      string `json:"phone"`
					NationalID string `json:"nationalID"`
					FullName   string `json:"fullName"`
				} `json:"contact"`
				Coordinates struct {
					Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Ending point latitude"`
					Longitude float64 `json:"longitude" example:"-70.6130425" description:"Ending point longitude"`
				} `json:"coordinates"`
				PoliticalArea struct {
					Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
					District string `json:"district" example:"la florida" description:"District name"`
					Province string `json:"province" example:"santiago" description:"Province name"`
					State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
				} `json:"politicalArea"`
				ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
			} `json:"addressInfo"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
		} `json:"endLocation"`
		Skills     []string `json:"skills" description:"Vehicle capabilities such as size or equipment requirements. eg: XL, heavy, etc"`
		TimeWindow struct {
			Start string `json:"start" example:"08:00" description:"Time window start (24h format)"`
			End   string `json:"end" example:"18:00" description:"Time window end (24h format)"`
		} `json:"timeWindow"`
		Capacity struct {
			Insurance             int64 `json:"insurance" example:"100000" description:"Maximum insurance value the vehicle can carry (CLP,MXN,PEN)"`
			Volume                int64 `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic centimeters (cm³)"`
			Weight                int64 `json:"weight" example:"1000" description:"Maximum weight in grams"`
			DeliveryUnitsQuantity int64 `json:"deliveryUnitsQuantity" example:"50" description:"Maximum number of delivery units the vehicle can carry"`
		} `json:"capacity"`
	} `json:"vehicles"`
	Visits []struct {
		Pickup struct {
			Instructions string `json:"instructions" example:"Recoger en recepción" description:"Instructions for pickup"`
			AddressInfo  struct {
				AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
				AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
				Contact      struct {
					Email      string `json:"email"`
					Phone      string `json:"phone"`
					NationalID string `json:"nationalID"`
					FullName   string `json:"fullName"`
				} `json:"contact"`
				Coordinates struct {
					Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Pickup point latitude"`
					Longitude float64 `json:"longitude" example:"-70.6130425" description:"Pickup point longitude"`
				} `json:"coordinates"`
				PoliticalArea struct {
					Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
					District string `json:"district" example:"la florida" description:"District name"`
					Province string `json:"province" example:"santiago" description:"Province name"`
					State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
				} `json:"politicalArea"`
				ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
			} `json:"addressInfo"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
			ServiceTime int64    `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
			Skills      []string `json:"skills" description:"Required vehicle capabilities for this visit"`
			TimeWindow  struct {
				Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
				End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
			} `json:"timeWindow"`
		} `json:"pickup"`
		Delivery struct {
			Instructions string `json:"instructions" example:"Entregar en recepción" description:"Instructions for delivery"`
			AddressInfo  struct {
				AddressLine1 string `json:"addressLine1" example:"Inglaterra 59" description:"Primary address line"`
				AddressLine2 string `json:"addressLine2" example:"Piso 2214" description:"Secondary address line"`
				Contact      struct {
					Email      string `json:"email"`
					Phone      string `json:"phone"`
					NationalID string `json:"nationalID"`
					FullName   string `json:"fullName"`
				} `json:"contact"`
				Coordinates struct {
					Latitude  float64 `json:"latitude" example:"-33.5147889" description:"Delivery point latitude"`
					Longitude float64 `json:"longitude" example:"-70.6130425" description:"Delivery point longitude"`
				} `json:"coordinates"`
				PoliticalArea struct {
					Code     string `json:"code" example:"cl-rm-la-florida" description:"Political area code"`
					District string `json:"district" example:"la florida" description:"District name"`
					Province string `json:"province" example:"santiago" description:"Province name"`
					State    string `json:"state" example:"region metropolitana de santiago" description:"State name"`
				} `json:"politicalArea"`
				ZipCode string `json:"zipCode" example:"7500000" description:"ZIP code"`
			} `json:"addressInfo"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
			ServiceTime int64    `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
			Skills      []string `json:"skills" description:"Required vehicle capabilities for this visit"`
			TimeWindow  struct {
				Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
				End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
			} `json:"timeWindow"`
		} `json:"delivery"`
		Orders []struct {
			DeliveryUnits []struct {
				Items []struct {
					Sku string `json:"sku" example:"SKU123" description:"Stock keeping unit identifier"`
				} `json:"items"`
				Insurance int64  `json:"insurance" example:"10000" description:"Insurance value of the delivery unit"`
				Volume    int64  `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic centimeters (cm³)"`
				Weight    int64  `json:"weight" example:"1000" description:"Weight of the delivery unit in grams"`
				Lpn       string `json:"lpn" example:"LPN456" description:"License plate number of the delivery unit"`
			} `json:"deliveryUnits"`
			ReferenceID string `json:"referenceID" example:"ORD789" description:"Unique identifier for the order"`
		} `json:"orders"`
	} `json:"visits"`
}

func (r *OptimizeFleetRequest) Map() optimization.FleetOptimization {
	vehicles := make([]optimization.Vehicle, len(r.Vehicles))
	for i, v := range r.Vehicles {
		vehicles[i] = optimization.Vehicle{
			Plate: v.Plate,
			StartLocation: optimization.AddressInfo{
				AddressLine1: v.StartLocation.AddressInfo.AddressLine1,
				AddressLine2: v.StartLocation.AddressInfo.AddressLine2,
				Contact: optimization.Contact{
					Email:      v.StartLocation.AddressInfo.Contact.Email,
					Phone:      v.StartLocation.AddressInfo.Contact.Phone,
					NationalID: v.StartLocation.AddressInfo.Contact.NationalID,
					FullName:   v.StartLocation.AddressInfo.Contact.FullName,
				},
				Coordinates: optimization.Coordinates{
					Latitude:  v.StartLocation.AddressInfo.Coordinates.Latitude,
					Longitude: v.StartLocation.AddressInfo.Coordinates.Longitude,
				},
				PoliticalArea: optimization.PoliticalArea{
					Code:     v.StartLocation.AddressInfo.PoliticalArea.Code,
					District: v.StartLocation.AddressInfo.PoliticalArea.District,
					Province: v.StartLocation.AddressInfo.PoliticalArea.Province,
					State:    v.StartLocation.AddressInfo.PoliticalArea.State,
				},
				ZipCode: v.StartLocation.AddressInfo.ZipCode,
			},
			EndLocation: optimization.AddressInfo{
				AddressLine1: v.EndLocation.AddressInfo.AddressLine1,
				AddressLine2: v.EndLocation.AddressInfo.AddressLine2,
				Contact: optimization.Contact{
					Email:      v.EndLocation.AddressInfo.Contact.Email,
					Phone:      v.EndLocation.AddressInfo.Contact.Phone,
					NationalID: v.EndLocation.AddressInfo.Contact.NationalID,
					FullName:   v.EndLocation.AddressInfo.Contact.FullName,
				},
				Coordinates: optimization.Coordinates{
					Latitude:  v.EndLocation.AddressInfo.Coordinates.Latitude,
					Longitude: v.EndLocation.AddressInfo.Coordinates.Longitude,
				},
				PoliticalArea: optimization.PoliticalArea{
					Code:     v.EndLocation.AddressInfo.PoliticalArea.Code,
					District: v.EndLocation.AddressInfo.PoliticalArea.District,
					Province: v.EndLocation.AddressInfo.PoliticalArea.Province,
					State:    v.EndLocation.AddressInfo.PoliticalArea.State,
				},
				ZipCode: v.EndLocation.AddressInfo.ZipCode,
			},
			Skills: v.Skills,
			TimeWindow: optimization.TimeWindow{
				Start: v.TimeWindow.Start,
				End:   v.TimeWindow.End,
			},
			Capacity: optimization.Capacity{
				Insurance:             v.Capacity.Insurance,
				Volume:                v.Capacity.Volume,
				Weight:                v.Capacity.Weight,
				DeliveryUnitsQuantity: v.Capacity.DeliveryUnitsQuantity,
			},
		}
	}

	visits := make([]optimization.Visit, len(r.Visits))
	for i, v := range r.Visits {
		// Mapear pickup
		pickup := optimization.VisitLocation{
			Instructions: v.Pickup.Instructions,
			AddressInfo: optimization.AddressInfo{
				AddressLine1: v.Pickup.AddressInfo.AddressLine1,
				AddressLine2: v.Pickup.AddressInfo.AddressLine2,
				Contact: optimization.Contact{
					Email:      v.Pickup.AddressInfo.Contact.Email,
					Phone:      v.Pickup.AddressInfo.Contact.Phone,
					NationalID: v.Pickup.AddressInfo.Contact.NationalID,
					FullName:   v.Pickup.AddressInfo.Contact.FullName,
				},
				Coordinates: optimization.Coordinates{
					Latitude:  v.Pickup.AddressInfo.Coordinates.Latitude,
					Longitude: v.Pickup.AddressInfo.Coordinates.Longitude,
				},
				PoliticalArea: optimization.PoliticalArea{
					Code:     v.Pickup.AddressInfo.PoliticalArea.Code,
					District: v.Pickup.AddressInfo.PoliticalArea.District,
					Province: v.Pickup.AddressInfo.PoliticalArea.Province,
					State:    v.Pickup.AddressInfo.PoliticalArea.State,
				},
				ZipCode: v.Pickup.AddressInfo.ZipCode,
			},
			NodeInfo: optimization.NodeInfo{
				ReferenceID: v.Pickup.NodeInfo.ReferenceID,
			},
			ServiceTime: v.Pickup.ServiceTime,
			Skills:      v.Pickup.Skills,
			TimeWindow: optimization.TimeWindow{
				Start: v.Pickup.TimeWindow.Start,
				End:   v.Pickup.TimeWindow.End,
			},
		}

		// Mapear delivery
		delivery := optimization.VisitLocation{
			Instructions: v.Delivery.Instructions,
			AddressInfo: optimization.AddressInfo{
				AddressLine1: v.Delivery.AddressInfo.AddressLine1,
				AddressLine2: v.Delivery.AddressInfo.AddressLine2,
				Contact: optimization.Contact{
					Email:      v.Delivery.AddressInfo.Contact.Email,
					Phone:      v.Delivery.AddressInfo.Contact.Phone,
					NationalID: v.Delivery.AddressInfo.Contact.NationalID,
					FullName:   v.Delivery.AddressInfo.Contact.FullName,
				},
				Coordinates: optimization.Coordinates{
					Latitude:  v.Delivery.AddressInfo.Coordinates.Latitude,
					Longitude: v.Delivery.AddressInfo.Coordinates.Longitude,
				},
				PoliticalArea: optimization.PoliticalArea{
					Code:     v.Delivery.AddressInfo.PoliticalArea.Code,
					District: v.Delivery.AddressInfo.PoliticalArea.District,
					Province: v.Delivery.AddressInfo.PoliticalArea.Province,
					State:    v.Delivery.AddressInfo.PoliticalArea.State,
				},
				ZipCode: v.Delivery.AddressInfo.ZipCode,
			},
			NodeInfo: optimization.NodeInfo{
				ReferenceID: v.Delivery.NodeInfo.ReferenceID,
			},
			ServiceTime: v.Delivery.ServiceTime,
			Skills:      v.Delivery.Skills,
			TimeWindow: optimization.TimeWindow{
				Start: v.Delivery.TimeWindow.Start,
				End:   v.Delivery.TimeWindow.End,
			},
		}

		// Mapear órdenes
		orders := make([]optimization.Order, len(v.Orders))
		for j, o := range v.Orders {
			deliveryUnits := make([]optimization.DeliveryUnit, len(o.DeliveryUnits))
			for k, du := range o.DeliveryUnits {
				items := make([]optimization.Item, len(du.Items))
				for l, item := range du.Items {
					items[l] = optimization.Item{
						Sku: item.Sku,
					}
				}
				deliveryUnits[k] = optimization.DeliveryUnit{
					Items:     items,
					Insurance: du.Insurance,
					Volume:    du.Volume,
					Weight:    du.Weight,
					Lpn:       du.Lpn,
				}
			}
			orders[j] = optimization.Order{
				DeliveryUnits: deliveryUnits,
				ReferenceID:   o.ReferenceID,
			}
		}

		visits[i] = optimization.Visit{
			Pickup:   pickup,
			Delivery: delivery,
			Orders:   orders,
		}
	}

	return optimization.FleetOptimization{
		PlanReferenceID: r.PlanReferenceID,
		Vehicles:        vehicles,
		Visits:          visits,
	}
}
