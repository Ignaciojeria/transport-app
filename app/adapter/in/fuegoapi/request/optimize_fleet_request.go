package request

import "transport-app/app/domain/optimization"

type OptimizeFleetRequest struct {
	PlanReferenceID string `json:"planReferenceID"`
	Vehicles        []struct {
		Plate         string `json:"plate" example:"SERV-80" description:"Vehicle license plate or internal code"`
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
		Skills     []string `json:"skills" description:"Vehicle capabilities such as size or equipment requirements. eg: XL, heavy, etc"`
		TimeWindow struct {
			Start string `json:"start" example:"08:00" description:"Time window start (24h format)"`
			End   string `json:"end" example:"18:00" description:"Time window end (24h format)"`
		} `json:"timeWindow"`
		Capacity struct {
			Insurance             int64 `json:"insurance" example:"100000" description:"Maximum insurance value the vehicle can carry (CLP,MXN,PEN)"`
			Volume                int64 `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
			Weight                int64 `json:"weight" example:"1000" description:"Maximum weight in grams"`
			DeliveryUnitsQuantity int64 `json:"deliveryUnitsQuantity" example:"50" description:"Maximum number of delivery units the vehicle can carry"`
		} `json:"capacity"`
	} `json:"vehicles"`
	Visits []struct {
		Pickup struct {
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
			Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
			TimeWindow struct {
				Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
				End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
			} `json:"timeWindow"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
		} `json:"pickup"`
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
			Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
			TimeWindow struct {
				Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
				End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
			} `json:"timeWindow"`
			NodeInfo struct {
				ReferenceID string `json:"referenceID"`
			} `json:"nodeInfo"`
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

func (r *OptimizeFleetRequest) Map() optimization.FleetOptimization {
	vehicles := make([]optimization.Vehicle, len(r.Vehicles))
	for i, v := range r.Vehicles {
		vehicles[i] = optimization.Vehicle{
			Plate: v.Plate,
			StartLocation: optimization.Location{
				Latitude:  v.StartLocation.Latitude,
				Longitude: v.StartLocation.Longitude,
				NodeInfo: optimization.NodeInfo{
					ReferenceID: v.StartLocation.NodeInfo.ReferenceID,
				},
			},
			EndLocation: optimization.Location{
				Latitude:  v.EndLocation.Latitude,
				Longitude: v.EndLocation.Longitude,
				NodeInfo: optimization.NodeInfo{
					ReferenceID: v.EndLocation.NodeInfo.ReferenceID,
				},
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
			Coordinates: optimization.Coordinates{
				Latitude:  v.Pickup.Coordinates.Latitude,
				Longitude: v.Pickup.Coordinates.Longitude,
			},
			ServiceTime: v.Pickup.ServiceTime,
			Contact: optimization.Contact{
				Email:      v.Pickup.Contact.Email,
				Phone:      v.Pickup.Contact.Phone,
				NationalID: v.Pickup.Contact.NationalID,
				FullName:   v.Pickup.Contact.FullName,
			},
			Skills: v.Pickup.Skills,
			TimeWindow: optimization.TimeWindow{
				Start: v.Pickup.TimeWindow.Start,
				End:   v.Pickup.TimeWindow.End,
			},
			NodeInfo: optimization.NodeInfo{
				ReferenceID: v.Pickup.NodeInfo.ReferenceID,
			},
		}

		// Mapear delivery
		delivery := optimization.VisitLocation{
			Coordinates: optimization.Coordinates{
				Latitude:  v.Delivery.Coordinates.Latitude,
				Longitude: v.Delivery.Coordinates.Longitude,
			},
			ServiceTime: v.Delivery.ServiceTime,
			Contact: optimization.Contact{
				Email:      v.Delivery.Contact.Email,
				Phone:      v.Delivery.Contact.Phone,
				NationalID: v.Delivery.Contact.NationalID,
				FullName:   v.Delivery.Contact.FullName,
			},
			Skills: v.Delivery.Skills,
			TimeWindow: optimization.TimeWindow{
				Start: v.Delivery.TimeWindow.Start,
				End:   v.Delivery.TimeWindow.End,
			},
			NodeInfo: optimization.NodeInfo{
				ReferenceID: v.Delivery.NodeInfo.ReferenceID,
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
