package model

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
	"transport-app/app/domain"
	"transport-app/app/domain/optimization"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/twpayne/go-polyline"
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

// decodeGeometry decodifica la geometría Polyline y retorna información sobre la ruta
func decodeGeometry(geometryStr string) (string, error) {
	if geometryStr == "" {
		return "Sin geometría disponible", nil
	}

	// Decodificar el polyline
	coords, _, err := polyline.DecodeCoords([]byte(geometryStr))
	if err != nil {
		return "", fmt.Errorf("error decodificando polyline: %w", err)
	}

	n := len(coords)
	if n == 0 {
		return "Polyline vacío", nil
	}

	// Construir la lista completa de puntos
	var pointsStr string
	for i, coord := range coords {
		if i > 0 {
			pointsStr += ", "
		}
		pointsStr += fmt.Sprintf("[%.6f, %.6f]", coord[1], coord[0]) // lon, lat
	}

	return fmt.Sprintf(
		"Polyline con %d puntos - Puntos completos: [%s]",
		n, pointsStr,
	), nil
}

func (ret VroomOptimizationResponse) Map(ctx context.Context, req optimization.FleetOptimization) domain.Plan {
	plan := domain.Plan{
		ReferenceID: uuid.New().String(),
		PlannedDate: time.Now(),
	}

	// Mapear rutas asignadas
	for _, vroomRoute := range ret.Routes {
		route := domain.Route{
			ReferenceID: uuid.New().String(),
		}

		// Decodificar y mostrar la geometría si está disponible
		if vroomRoute.Geometry != "" {
			geometryInfo, err := decodeGeometry(vroomRoute.Geometry)
			if err != nil {
				fmt.Printf("Error decodificando geometría para ruta %d: %v\n", vroomRoute.Vehicle, err)
			} else {
				fmt.Printf("Geometría de ruta %d: %s\n", vroomRoute.Vehicle, geometryInfo)
			}
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

	if err := ret.ExportToGeoJSON("ui/static/dev/geojson.json"); err != nil {
		fmt.Printf("error exportando GeoJSON: %v\n", err)
	}

	return plan
}

// findVisitByJobID busca una visita que corresponde a un job (solo delivery válido)
func findVisitByJobID(jobID int64, visits []optimization.Visit) *optimization.Visit {
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
func findVisitByShipmentID(shipmentID int64, visits []optimization.Visit) *optimization.Visit {
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
func createOrdersFromVisit(visit *optimization.Visit, hasPickup bool) []domain.Order {
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

// GeoJSONFeature representa una feature de GeoJSON
type GeoJSONFeature struct {
	Type       string                 `json:"type"`
	Geometry   GeoJSONGeometry        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// GeoJSONGeometry representa la geometría de GeoJSON
type GeoJSONGeometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

// GeoJSONCollection representa una colección de features de GeoJSON
type GeoJSONCollection struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}

// ExportToGeoJSON exporta las rutas a un archivo GeoJSON
func (ret VroomOptimizationResponse) ExportToGeoJSON(filename string) error {
	collection := GeoJSONCollection{
		Type:     "FeatureCollection",
		Features: []GeoJSONFeature{},
	}

	// Procesar cada ruta
	for i, route := range ret.Routes {
		if route.Geometry == "" {
			continue
		}

		// Decodificar el polyline
		coords, _, err := polyline.DecodeCoords([]byte(route.Geometry))
		if err != nil {
			fmt.Printf("Error decodificando polyline para ruta %d: %v\n", i, err)
			continue
		}

		// Convertir coordenadas a formato GeoJSON (lon, lat)
		var geoJSONCoords [][]float64
		for _, coord := range coords {
			geoJSONCoords = append(geoJSONCoords, []float64{coord[1], coord[0]}) // lon, lat
		}

		// Crear feature para la ruta
		routeFeature := GeoJSONFeature{
			Type: "Feature",
			Geometry: GeoJSONGeometry{
				Type:        "LineString",
				Coordinates: geoJSONCoords,
			},
			Properties: map[string]interface{}{
				"vehicle":        route.Vehicle,
				"cost":           route.Cost,
				"service":        route.Service,
				"duration":       route.Duration,
				"waiting_time":   route.WaitingTime,
				"priority":       route.Priority,
				"steps_count":    len(route.Steps),
				"name":           fmt.Sprintf("Ruta del vehículo %d", route.Vehicle),
				"stroke":         "#FF0000", // Color de la línea
				"stroke-width":   3,         // Grosor de la línea
				"stroke-opacity": 0.8,       // Opacidad de la línea
			},
		}

		collection.Features = append(collection.Features, routeFeature)

		// Contadores para secuencias independientes por tipo
		pickupCount := 1
		deliveryCount := 1
		jobCount := 1

		// Agregar puntos para cada step con ubicación y símbolos de secuencia
		for j, step := range route.Steps {
			if len(step.Location) == 2 {
				// Determinar el símbolo y color basado en el tipo de paso con secuencia independiente
				symbol, color, size := getStepSymbolWithSequence(step.Type, j, pickupCount, deliveryCount, jobCount)

				// Incrementar el contador correspondiente
				switch step.Type {
				case "pickup":
					pickupCount++
				case "delivery":
					deliveryCount++
				case "job":
					jobCount++
				}

				stepFeature := GeoJSONFeature{
					Type: "Feature",
					Geometry: GeoJSONGeometry{
						Type:        "Point",
						Coordinates: []float64{step.Location[0], step.Location[1]}, // lon, lat
					},
					Properties: map[string]interface{}{
						"vehicle":     route.Vehicle,
						"step_index":  j,
						"step_number": j,
						"step_type":   step.Type,
						"arrival":     step.Arrival,
						"duration":    step.Duration,
						"service":     step.Service,
						"job":         step.Job,
						"shipment":    step.Shipment,
						"description": step.Description,
						"name":        fmt.Sprintf("Paso %d - %s", j, step.Type),
						// Propiedades para símbolos
						"marker-symbol": symbol,
						"marker-color":  color,
						"marker-size":   size,
						"title":         fmt.Sprintf("Vehículo %d - Paso %d: %s", route.Vehicle, j, step.Type),
						"popup":         fmt.Sprintf("Vehículo: %d<br>Paso: %d<br>Tipo: %s<br>Llegada: %d seg", route.Vehicle, j, step.Type, step.Arrival),
					},
				}

				collection.Features = append(collection.Features, stepFeature)
			}
		}
	}

	// Agregar puntos para trabajos no asignados
	for _, unassigned := range ret.Unassigned {
		if len(unassigned.Location) == 2 {
			unassignedFeature := GeoJSONFeature{
				Type: "Feature",
				Geometry: GeoJSONGeometry{
					Type:        "Point",
					Coordinates: []float64{unassigned.Location[0], unassigned.Location[1]}, // lon, lat
				},
				Properties: map[string]interface{}{
					"job_id": unassigned.ID,
					"reason": unassigned.Reason,
					"status": "unassigned",
					"name":   fmt.Sprintf("Trabajo no asignado %d", unassigned.ID),
					// Propiedades para símbolos de trabajos no asignados
					"marker-symbol": "cross",
					"marker-color":  "#FF0000",
					"marker-size":   "large",
					"title":         fmt.Sprintf("Trabajo no asignado %d", unassigned.ID),
					"popup":         fmt.Sprintf("Trabajo: %d<br>Razón: %s", unassigned.ID, unassigned.Reason),
				},
			}

			collection.Features = append(collection.Features, unassignedFeature)
		}
	}

	// Convertir a JSON
	jsonData, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando GeoJSON: %w", err)
	}

	// Asegurarse de que el directorio de destino exista
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creando directorio %s: %w", dir, err)
	}

	// Escribir al archivo
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo GeoJSON: %w", err)
	}

	fmt.Printf("GeoJSON exportado exitosamente a: %s\n", filename)
	fmt.Printf("Total de features: %d\n", len(collection.Features))

	return nil
}

// getStepSymbolWithSequence determina el símbolo, color y tamaño basado en el tipo de paso con secuencias independientes
func getStepSymbolWithSequence(stepType string, stepNumber int, pickupCount, deliveryCount, jobCount int) (symbol, color, size string) {
	switch stepType {
	case "start":
		return "play", "#00FF00", "large" // Verde para inicio
	case "end":
		return "stop", "#FF0000", "large" // Rojo para fin
	case "pickup":
		return fmt.Sprintf("%d", pickupCount), "#0000FF", "medium" // Azul con número secuencial de pickup
	case "delivery":
		return fmt.Sprintf("%d", deliveryCount), "#FF6600", "medium" // Naranja con número secuencial de delivery
	case "job":
		return fmt.Sprintf("%d", jobCount), "#9932CC", "medium" // Púrpura con número secuencial de job
	case "break":
		return "coffee", "#8B4513", "medium" // Marrón para descanso
	default:
		return fmt.Sprintf("%d", stepNumber+1), "#666666", "medium" // Gris por defecto con número secuencial
	}
}
