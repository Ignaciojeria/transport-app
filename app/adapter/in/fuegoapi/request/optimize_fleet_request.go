package request

import "transport-app/app/domain/optimization"

// Estructuras granulares para OptimizeFleetRequest

type OptimizeFleetRequest struct {
	PlanReferenceID string                 `json:"planReferenceID"`
	Vehicles        []OptimizeFleetVehicle `json:"vehicles"`
	Visits          []OptimizeFleetVisit   `json:"visits"`
}

type OptimizeFleetVehicle struct {
	Plate         string                       `json:"plate"`
	StartLocation OptimizeFleetVehicleLocation `json:"startLocation"`
	EndLocation   OptimizeFleetVehicleLocation `json:"endLocation"`
	Skills        []string                     `json:"skills"`
	TimeWindow    OptimizeFleetTimeWindow      `json:"timeWindow"`
	Capacity      OptimizeFleetVehicleCapacity `json:"capacity"`
}

type OptimizeFleetVehicleLocation struct {
	AddressInfo OptimizeFleetAddressInfo `json:"addressInfo"`
	NodeInfo    OptimizeFleetNodeInfo    `json:"nodeInfo"`
}

type OptimizeFleetAddressInfo struct {
	AddressLine1  string                     `json:"addressLine1"`
	AddressLine2  string                     `json:"addressLine2"`
	Contact       OptimizeFleetContact       `json:"contact"`
	Coordinates   OptimizeFleetCoordinates   `json:"coordinates"`
	PoliticalArea OptimizeFleetPoliticalArea `json:"politicalArea"`
	ZipCode       string                     `json:"zipCode"`
}

type OptimizeFleetContact struct {
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	NationalID string `json:"nationalID"`
	FullName   string `json:"fullName"`
}

type OptimizeFleetCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type OptimizeFleetPoliticalArea struct {
	Code            string `json:"code"`
	AdminAreaLevel1 string `json:"adminAreaLevel1"`
	AdminAreaLevel2 string `json:"adminAreaLevel2"`
	AdminAreaLevel3 string `json:"adminAreaLevel3"`
	AdminAreaLevel4 string `json:"adminAreaLevel4"`
}

type OptimizeFleetVehicleCapacity struct {
	Volume                int64 `json:"volume" example:"1000" description:"Volume in cubic centimeters (cm³)"`
	Weight                int64 `json:"weight" example:"1000" description:"Weight in grams (g)"`
	Insurance             int64 `json:"insurance" example:"10000" description:"Insurance value in currency units (CLP, MXN, PEN, CENTS etc.) - only integer values accepted"`
	DeliveryUnitsQuantity int64 `json:"deliveryUnitsQuantity"`
}

type OptimizeFleetTimeWindow struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type OptimizeFleetVisit struct {
	Pickup   OptimizeFleetVisitLocation `json:"pickup"`
	Delivery OptimizeFleetVisitLocation `json:"delivery"`
	Orders   []OptimizeFleetOrder       `json:"orders"`
}

type OptimizeFleetVisitLocation struct {
	Instructions string                   `json:"instructions"`
	AddressInfo  OptimizeFleetAddressInfo `json:"addressInfo"`
	NodeInfo     OptimizeFleetNodeInfo    `json:"nodeInfo"`
	ServiceTime  int64                    `json:"serviceTime"`
	TimeWindow   OptimizeFleetTimeWindow  `json:"timeWindow"`
}

type OptimizeFleetNodeInfo struct {
	ReferenceID string `json:"referenceID"`
}

type OptimizeFleetOrder struct {
	DeliveryUnits []OptimizeFleetDeliveryUnit `json:"deliveryUnits"`
	ReferenceID   string                      `json:"referenceID"`
}

type OptimizeFleetDeliveryUnit struct {
	Items     []OptimizeFleetItem `json:"items"`
	Volume    int64               `json:"volume" example:"1000" description:"Volume in cubic centimeters (cm³)"`
	Weight    int64               `json:"weight" example:"1000" description:"Weight in grams (g)"`
	Insurance int64               `json:"insurance" example:"10000" description:"Insurance value in currency units (CLP, MXN, PEN, CENTS etc.) - only integer values accepted"`
	Lpn       string              `json:"lpn"`
	Skills    []string            `json:"skills"`
}

type OptimizeFleetItem struct {
	Sku string `json:"sku"`
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
					Code:            v.StartLocation.AddressInfo.PoliticalArea.Code,
					AdminAreaLevel1: v.StartLocation.AddressInfo.PoliticalArea.AdminAreaLevel1,
					AdminAreaLevel2: v.StartLocation.AddressInfo.PoliticalArea.AdminAreaLevel2,
					AdminAreaLevel3: v.StartLocation.AddressInfo.PoliticalArea.AdminAreaLevel3,
					AdminAreaLevel4: v.StartLocation.AddressInfo.PoliticalArea.AdminAreaLevel4,
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
					Code:            v.EndLocation.AddressInfo.PoliticalArea.Code,
					AdminAreaLevel1: v.EndLocation.AddressInfo.PoliticalArea.AdminAreaLevel1,
					AdminAreaLevel2: v.EndLocation.AddressInfo.PoliticalArea.AdminAreaLevel2,
					AdminAreaLevel3: v.EndLocation.AddressInfo.PoliticalArea.AdminAreaLevel3,
					AdminAreaLevel4: v.EndLocation.AddressInfo.PoliticalArea.AdminAreaLevel4,
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
					Code:            v.Pickup.AddressInfo.PoliticalArea.Code,
					AdminAreaLevel1: v.Pickup.AddressInfo.PoliticalArea.AdminAreaLevel1,
					AdminAreaLevel2: v.Pickup.AddressInfo.PoliticalArea.AdminAreaLevel2,
					AdminAreaLevel3: v.Pickup.AddressInfo.PoliticalArea.AdminAreaLevel3,
					AdminAreaLevel4: v.Pickup.AddressInfo.PoliticalArea.AdminAreaLevel4,
				},
				ZipCode: v.Pickup.AddressInfo.ZipCode,
			},
			NodeInfo: optimization.NodeInfo{
				ReferenceID: v.Pickup.NodeInfo.ReferenceID,
			},
			ServiceTime: v.Pickup.ServiceTime,
			TimeWindow: optimization.TimeWindow{
				Start: v.Pickup.TimeWindow.Start,
				End:   v.Pickup.TimeWindow.End,
			},
		}

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
					Code:            v.Delivery.AddressInfo.PoliticalArea.Code,
					AdminAreaLevel1: v.Delivery.AddressInfo.PoliticalArea.AdminAreaLevel1,
					AdminAreaLevel2: v.Delivery.AddressInfo.PoliticalArea.AdminAreaLevel2,
					AdminAreaLevel3: v.Delivery.AddressInfo.PoliticalArea.AdminAreaLevel3,
					AdminAreaLevel4: v.Delivery.AddressInfo.PoliticalArea.AdminAreaLevel4,
				},
				ZipCode: v.Delivery.AddressInfo.ZipCode,
			},
			NodeInfo: optimization.NodeInfo{
				ReferenceID: v.Delivery.NodeInfo.ReferenceID,
			},
			ServiceTime: v.Delivery.ServiceTime,
			TimeWindow: optimization.TimeWindow{
				Start: v.Delivery.TimeWindow.Start,
				End:   v.Delivery.TimeWindow.End,
			},
		}

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
					Skills:    du.Skills,
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
