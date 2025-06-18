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
					UnitOfMeasure: "kg",
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

		// Mapear órdenes basadas en los jobs de los steps
		var orders []domain.Order
		for _, step := range vroomRoute.Steps {
			if step.Job > 0 && step.Job <= int64(len(req.Visits)) {
				visit := req.Visits[step.Job-1] // VROOM usa índices basados en 1

				// Crear orden para cada order en la visita
				for _, orderReq := range visit.Orders {
					order := domain.Order{
						ReferenceID: domain.ReferenceID(orderReq.ReferenceID),
						Origin: domain.NodeInfo{
							ReferenceID: domain.ReferenceID(uuid.New().String()),
							Name:        "Origen de recogida",
							AddressInfo: domain.AddressInfo{
								Coordinates: domain.Coordinates{
									Point: orb.Point{visit.PickupLocation.Longitude, visit.PickupLocation.Latitude},
								},
							},
						},
						Destination: domain.NodeInfo{
							ReferenceID: domain.ReferenceID(uuid.New().String()),
							Name:        "Destino de entrega",
							AddressInfo: domain.AddressInfo{
								Coordinates: domain.Coordinates{
									Point: orb.Point{visit.DispatchLocation.Longitude, visit.DispatchLocation.Latitude},
								},
							},
						},
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
			}
		}
		route.Orders = orders

		plan.Routes = append(plan.Routes, route)
	}

	// Mapear trabajos no asignados
	for _, unassigned := range ret.Unassigned {
		if unassigned.ID > 0 && unassigned.ID <= int64(len(req.Visits)) {
			visit := req.Visits[unassigned.ID-1]

			// Crear órdenes no asignadas
			for _, orderReq := range visit.Orders {
				order := domain.Order{
					ReferenceID:      domain.ReferenceID(orderReq.ReferenceID),
					UnassignedReason: unassigned.Reason,
					Origin: domain.NodeInfo{
						ReferenceID: domain.ReferenceID(uuid.New().String()),
						Name:        "Origen de recogida",
						AddressInfo: domain.AddressInfo{
							Coordinates: domain.Coordinates{
								Point: orb.Point{visit.PickupLocation.Longitude, visit.PickupLocation.Latitude},
							},
						},
					},
					Destination: domain.NodeInfo{
						ReferenceID: domain.ReferenceID(uuid.New().String()),
						Name:        "Destino de entrega",
						AddressInfo: domain.AddressInfo{
							Coordinates: domain.Coordinates{
								Point: orb.Point{visit.DispatchLocation.Longitude, visit.DispatchLocation.Latitude},
							},
						},
					},
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

				plan.UnassignedOrders = append(plan.UnassignedOrders, order)
			}
		}
	}

	return plan
}
