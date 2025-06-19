package model

import (
	"context"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/domain"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

// VroomOptimizationResponse represents the complete response from VROOM API
type VroomOptimizationResponse struct {
	Code       int64           `json:"code"`
	Error      string          `json:"error,omitempty"`
	Unassigned []UnassignedJob `json:"unassigned,omitempty"`
	Routes     []Route         `json:"routes,omitempty"`
}

// UnassignedJob represents jobs that couldn't be assigned to any vehicle
type UnassignedJob struct {
	ID       int64      `json:"id"`
	Location [2]float64 `json:"location,omitempty"`
	Reason   string     `json:"reason"`
}

// Route represents a vehicle's route with its assigned jobs
type Route struct {
	Vehicle     int64   `json:"vehicle"`
	Cost        int64   `json:"cost"`
	Service     int64   `json:"service"`
	Duration    int64   `json:"duration"`
	WaitingTime int64   `json:"waiting_time"`
	Priority    float64 `json:"priority"`
	Steps       []Step  `json:"steps"`
	Geometry    string  `json:"geometry,omitempty"`
}

// Step represents a single step in a route (pickup, delivery, or break)
type Step struct {
	Type           string         `json:"type"` // "start", "job", "pickup", "delivery", "break", "end"
	Arrival        int64          `json:"arrival"`
	Duration       int64          `json:"duration"`
	Service        int64          `json:"service"`
	WaitingTime    int64          `json:"waiting_time"`
	Job            int64          `json:"job,omitempty"`
	Location       [2]float64     `json:"location,omitempty"`
	Load           []int64        `json:"load,omitempty"`
	Distance       int64          `json:"distance,omitempty"`
	Setup          int64          `json:"setup,omitempty"`
	Shipment       int64          `json:"shipment,omitempty"`
	Pickup         int64          `json:"pickup,omitempty"`
	Delivery       int64          `json:"delivery,omitempty"`
	Description    string         `json:"description,omitempty"`
	CustomUserData map[string]any `json:"custom_user_data,omitempty"`
}

func (ret VroomOptimizationResponse) Map(ctx context.Context, req request.OptimizationRequest) domain.Plan {
	plan := domain.Plan{
		ReferenceID: uuid.New().String(),
		PlannedDate: time.Now(),
	}

	// Mapear rutas asignadas
	for _, vroomRoute := range ret.Routes {
		route := domain.Route{
			ReferenceID: uuid.New().String(),
		}

		// Mapear vehículo si existe en la solicitud
		if vroomRoute.Vehicle < int64(len(req.Vehicles)) {
			vehicle := req.Vehicles[vroomRoute.Vehicle]
			route.Vehicle = domain.Vehicle{
				Plate: vehicle.Plate,
				Weight: struct {
					Value         int
					UnitOfMeasure string
				}{
					Value:         int(vehicle.Capacity.Weight),
					UnitOfMeasure: "g", // El modelo usa gramos
				},
				Insurance: struct {
					PolicyStartDate      string
					PolicyExpirationDate string
					PolicyRenewalDate    string
					MaxInsuranceCoverage struct {
						Amount   float64
						Currency string
					}
				}{
					MaxInsuranceCoverage: struct {
						Amount   float64
						Currency string
					}{
						Amount:   float64(vehicle.Capacity.Insurance),
						Currency: "CLP", // Asumiendo pesos chilenos por defecto
					},
				},
			}
		}

		// Mapear origen y destino basado en los steps
		if len(vroomRoute.Steps) > 0 {
			// Origen: primer step con location
			for _, step := range vroomRoute.Steps {
				if step.Type == "start" && len(step.Location) == 2 {
					route.Origin = domain.NodeInfo{
						ReferenceID: domain.ReferenceID(uuid.New().String()),
						Name:        "Origen",
						AddressInfo: domain.AddressInfo{
							Coordinates: domain.Coordinates{
								Point: orb.Point{step.Location[0], step.Location[1]},
							},
						},
					}
					break
				}
			}

			// Destino: último step con location
			for i := len(vroomRoute.Steps) - 1; i >= 0; i-- {
				step := vroomRoute.Steps[i]
				if step.Type == "end" && len(step.Location) == 2 {
					route.Destination = domain.NodeInfo{
						ReferenceID: domain.ReferenceID(uuid.New().String()),
						Name:        "Destino",
						AddressInfo: domain.AddressInfo{
							Coordinates: domain.Coordinates{
								Point: orb.Point{step.Location[0], step.Location[1]},
							},
						},
					}
					break
				}
			}
		}

		// Mapear órdenes basadas en los jobs y shipments de los steps
		var orders []domain.Order
		for _, step := range vroomRoute.Steps {
			// Manejar jobs (solo delivery)
			if step.Job > 0 {
				visit := findVisitByJobID(step.Job, req.Visits)
				if visit != nil {
					orders = append(orders, createOrdersFromVisit(visit, false)...)
				}
			}

			// Manejar shipments (pickup y delivery)
			if step.Shipment > 0 {
				visit := findVisitByShipmentID(step.Shipment, req.Visits)
				if visit != nil {
					orders = append(orders, createOrdersFromVisit(visit, true)...)
				}
			}
		}
		route.Orders = orders

		plan.Routes = append(plan.Routes, route)
	}

	// Mapear trabajos no asignados
	for _, unassigned := range ret.Unassigned {
		visit := findVisitByJobID(unassigned.ID, req.Visits)
		if visit != nil {
			unassignedOrders := createOrdersFromVisit(visit, false)
			// Marcar como no asignadas
			for i := range unassignedOrders {
				unassignedOrders[i].UnassignedReason = unassigned.Reason
			}
			plan.UnassignedOrders = append(plan.UnassignedOrders, unassignedOrders...)
		}
	}

	return plan
}

// findVisitByJobID busca una visita que corresponde a un job (solo delivery válido)
func findVisitByJobID(jobID int64, visits []struct {
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
}) *struct {
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
} {
	// Buscar la visita que corresponde a este job ID
	for i, v := range visits {
		// Verificar si esta visita tiene solo delivery válido (job)
		hasValidPickup := v.Pickup.Coordinates.Longitude != 0 || v.Pickup.Coordinates.Latitude != 0
		hasValidDelivery := v.Delivery.Coordinates.Longitude != 0 || v.Delivery.Coordinates.Latitude != 0

		if !hasValidPickup && hasValidDelivery {
			// Esta visita corresponde a un job
			return &visits[i]
		}
	}
	return nil
}

// findVisitByShipmentID busca una visita que corresponde a un shipment (pickup y delivery válidos)
func findVisitByShipmentID(shipmentID int64, visits []struct {
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
}) *struct {
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
} {
	// Buscar la visita que corresponde a este shipment ID
	for i, v := range visits {
		// Verificar si esta visita tiene pickup y delivery válidos (shipment)
		hasValidPickup := v.Pickup.Coordinates.Longitude != 0 || v.Pickup.Coordinates.Latitude != 0
		hasValidDelivery := v.Delivery.Coordinates.Longitude != 0 || v.Delivery.Coordinates.Latitude != 0

		if hasValidPickup && hasValidDelivery {
			// Esta visita corresponde a un shipment
			return &visits[i]
		}
	}
	return nil
}

// createOrdersFromVisit crea órdenes del dominio basadas en una visita
func createOrdersFromVisit(visit *struct {
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
}, hasPickup bool) []domain.Order {
	var orders []domain.Order

	// Crear orden para cada order en la visita
	for _, orderReq := range visit.Orders {
		order := domain.Order{
			ReferenceID: domain.ReferenceID(orderReq.ReferenceID),
			Destination: domain.NodeInfo{
				ReferenceID: domain.ReferenceID(uuid.New().String()),
				Name:        "Destino de entrega",
				AddressInfo: domain.AddressInfo{
					Coordinates: domain.Coordinates{
						Point: orb.Point{visit.Delivery.Coordinates.Longitude, visit.Delivery.Coordinates.Latitude},
					},
				},
			},
		}

		// Para shipments (pickup + delivery), incluir origen
		if hasPickup {
			order.Origin = domain.NodeInfo{
				ReferenceID: domain.ReferenceID(uuid.New().String()),
				Name:        "Origen de recogida",
				AddressInfo: domain.AddressInfo{
					Coordinates: domain.Coordinates{
						Point: orb.Point{visit.Pickup.Coordinates.Longitude, visit.Pickup.Coordinates.Latitude},
					},
				},
			}
		}

		// Mapear delivery units
		var deliveryUnits domain.DeliveryUnits
		for _, duReq := range orderReq.DeliveryUnits {
			deliveryUnit := domain.DeliveryUnit{
				Lpn: duReq.Lpn,
			}

			// Mapear items
			for _, itemReq := range duReq.Items {
				item := domain.Item{
					Sku: itemReq.Sku,
				}
				deliveryUnit.Items = append(deliveryUnit.Items, item)
			}

			deliveryUnits = append(deliveryUnits, deliveryUnit)
		}
		order.DeliveryUnits = deliveryUnits

		orders = append(orders, order)
	}

	return orders
}
