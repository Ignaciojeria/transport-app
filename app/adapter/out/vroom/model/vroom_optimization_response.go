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
			_, err := decodeGeometry(vroomRoute.Geometry)
			if err != nil {
				fmt.Printf("Error decodificando geometría para ruta %d: %v\n", vroomRoute.Vehicle, err)
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
		// Primero intentar encontrar como job
		visit := findVisitByJobID(unassigned.ID, req.Visits)
		if visit != nil {
			unassignedOrders := createOrdersFromVisit(visit, false)
			// Marcar como no asignadas
			for i := range unassignedOrders {
				unassignedOrders[i].UnassignedReason = unassigned.Reason
			}
			plan.UnassignedOrders = append(plan.UnassignedOrders, unassignedOrders...)
		} else {
			// Si no se encuentra como job, intentar como shipment
			visit = findVisitByShipmentID(unassigned.ID, req.Visits)
			if visit != nil {
				unassignedOrders := createOrdersFromVisit(visit, true)
				// Marcar como no asignadas
				for i := range unassignedOrders {
					unassignedOrders[i].UnassignedReason = unassigned.Reason
				}
				plan.UnassignedOrders = append(plan.UnassignedOrders, unassignedOrders...)
			} else {
				// Si no se encuentra ni como job ni como shipment, loggear el problema
				fmt.Printf("No se pudo encontrar la visita para el trabajo no asignado ID %d\n", unassigned.ID)
			}
		}
	}

	if err := ret.ExportToPolylineJSON("ui/static/dev/polyline.json", req); err != nil {
		fmt.Printf("error exportando datos de ruta: %v\n", err)
	}

	return plan
}

// findVisitByJobID busca una visita que corresponde a un job (solo delivery válido)
func findVisitByJobID(jobID int64, visits []optimization.Visit) *optimization.Visit {
	// Los job IDs en VROOM corresponden al índice de la visita en el request original
	// jobID 1 = primera visita, jobID 2 = segunda visita, etc.
	// Pero necesitamos ajustar porque los jobs se crean solo para visitas sin pickup válido

	jobIndex := 0
	for i, v := range visits {
		// Verificar si esta visita tiene solo delivery válido (job)
		hasValidPickup := v.Pickup.Coordinates.Longitude != 0 || v.Pickup.Coordinates.Latitude != 0
		hasValidDelivery := v.Delivery.Coordinates.Longitude != 0 || v.Delivery.Coordinates.Latitude != 0

		if !hasValidPickup && hasValidDelivery {
			// Esta visita corresponde a un job
			jobIndex++
			if int64(jobIndex) == jobID {
				return &visits[i]
			}
		}
	}
	return nil
}

// findVisitByShipmentID busca una visita que corresponde a un shipment (pickup y delivery válidos)
func findVisitByShipmentID(shipmentID int64, visits []optimization.Visit) *optimization.Visit {
	// Los shipment IDs en VROOM corresponden al índice de la visita en el request original
	// shipmentID 1 = primera visita, shipmentID 2 = segunda visita, etc.
	// Pero necesitamos ajustar porque los shipments se crean solo para visitas con pickup válido

	shipmentIndex := 0
	for i, v := range visits {
		// Verificar si esta visita tiene pickup y delivery válidos (shipment)
		hasValidPickup := v.Pickup.Coordinates.Longitude != 0 || v.Pickup.Coordinates.Latitude != 0
		hasValidDelivery := v.Delivery.Coordinates.Longitude != 0 || v.Delivery.Coordinates.Latitude != 0

		if hasValidPickup && hasValidDelivery {
			// Esta visita corresponde a un shipment
			shipmentIndex++
			if int64(shipmentIndex) == shipmentID {
				return &visits[i]
			}
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
func (ret VroomOptimizationResponse) ExportToGeoJSON(filename string, req optimization.FleetOptimization) error {
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

		// Agregar puntos para cada step con ubicación usando secuencia completa de la ruta
		stepCounter := 1 // Contador que empieza en 1 para steps con ubicación
		for _, step := range route.Steps {
			if len(step.Location) == 2 {
				symbol, color, size := getStepSymbolWithSequence(step.Type, stepCounter)

				// Extraer referenceID de las órdenes
				orderRefs := extractOrderRefsFromStep(step, req)

				stepFeature := GeoJSONFeature{
					Type: "Feature",
					Geometry: GeoJSONGeometry{
						Type:        "Point",
						Coordinates: []float64{step.Location[0], step.Location[1]}, // lon, lat
					},
					Properties: map[string]interface{}{
						"vehicle":     route.Vehicle,
						"step_index":  stepCounter - 1,
						"step_number": stepCounter, // Usar el contador separado
						"step_type":   step.Type,
						"arrival":     step.Arrival,
						"duration":    step.Duration,
						"service":     step.Service,
						"job":         step.Job,
						"shipment":    step.Shipment,
						"pickup":      step.Pickup,
						"delivery":    step.Delivery,
						"description": step.Description,
						"order_refs":  orderRefs, // ReferenceIDs de las órdenes
						"name":        fmt.Sprintf("Paso %d - %s", stepCounter, step.Type),
						// Propiedades para símbolos
						"marker-symbol": symbol,
						"marker-color":  color,
						"marker-size":   size,
						"title":         fmt.Sprintf("Vehículo %d - Paso %d: %s", route.Vehicle, stepCounter, step.Type),
						"popup":         fmt.Sprintf("Vehículo: %d<br>Paso: %d<br>Tipo: %s<br>Llegada: %d seg<br>Órdenes: %v", route.Vehicle, stepCounter, step.Type, step.Arrival, orderRefs),
					},
				}

				collection.Features = append(collection.Features, stepFeature)
				stepCounter++ // Incrementar solo cuando se procesa un step con ubicación
			}
		}
	}

	// Agregar puntos para trabajos no asignados
	for _, unassigned := range ret.Unassigned {
		if len(unassigned.Location) == 2 {
			// Buscar órdenes asociadas al trabajo no asignado
			var orderRefs []string
			visit := findVisitByJobID(unassigned.ID, req.Visits)
			if visit != nil {
				for _, order := range visit.Orders {
					orderRefs = append(orderRefs, order.ReferenceID)
				}
			} else {
				// Si no se encuentra como job, intentar como shipment
				visit = findVisitByShipmentID(unassigned.ID, req.Visits)
				if visit != nil {
					for _, order := range visit.Orders {
						orderRefs = append(orderRefs, order.ReferenceID)
					}
				}
			}

			unassignedFeature := GeoJSONFeature{
				Type: "Feature",
				Geometry: GeoJSONGeometry{
					Type:        "Point",
					Coordinates: []float64{unassigned.Location[0], unassigned.Location[1]}, // lon, lat
				},
				Properties: map[string]interface{}{
					"job_id":     unassigned.ID,
					"reason":     unassigned.Reason,
					"status":     "unassigned",
					"order_refs": orderRefs, // ReferenceIDs de las órdenes no asignadas
					"name":       fmt.Sprintf("Trabajo no asignado %d", unassigned.ID),
					// Propiedades para símbolos de trabajos no asignados
					"marker-symbol": "cross",
					"marker-color":  "#FF0000",
					"marker-size":   "large",
					"title":         fmt.Sprintf("Trabajo no asignado %d", unassigned.ID),
					"popup":         fmt.Sprintf("Trabajo: %d<br>Razón: %s<br>Órdenes: %v", unassigned.ID, unassigned.Reason, orderRefs),
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

// getStepSymbolWithSequence determina el símbolo, color y tamaño basado en el tipo de paso con secuencia completa
func getStepSymbolWithSequence(stepType string, stepNumber int) (symbol, color, size string) {
	switch stepType {
	case "start":
		return "play", "#00FF00", "large" // Verde para inicio
	case "end":
		return "stop", "#FF0000", "large" // Rojo para fin
	case "pickup":
		return fmt.Sprintf("%d", stepNumber), "#0000FF", "medium" // Azul con número secuencial completo
	case "delivery":
		return fmt.Sprintf("%d", stepNumber), "#FF6600", "medium" // Naranja con número secuencial completo
	case "job":
		return fmt.Sprintf("%d", stepNumber), "#9932CC", "medium" // Púrpura con número secuencial completo
	case "break":
		return "coffee", "#8B4513", "medium" // Marrón para descanso
	default:
		return fmt.Sprintf("%d", stepNumber), "#666666", "medium" // Gris por defecto con número secuencial completo
	}
}

// RouteData representa los datos de ruta para el frontend
type RouteData struct {
	Route      [][]float64       `json:"route"` // Coordenadas decodificadas del polyline
	Steps      []StepPoint       `json:"steps"` // Puntos de parada
	Vehicle    int64             `json:"vehicle"`
	Cost       int64             `json:"cost"`
	Duration   int64             `json:"duration"`
	Unassigned []UnassignedPoint `json:"unassigned,omitempty"`
}

// StepPoint representa un punto de parada
type StepPoint struct {
	Location    [2]float64 `json:"location"`
	StepType    string     `json:"step_type"`
	StepNumber  int        `json:"step_number"`
	Arrival     int64      `json:"arrival"`
	Description string     `json:"description,omitempty"`
	OrderRefs   []string   `json:"order_refs,omitempty"`  // ReferenceIDs de las órdenes asociadas
	JobID       int64      `json:"job_id,omitempty"`      // ID del job si aplica
	ShipmentID  int64      `json:"shipment_id,omitempty"` // ID del shipment si aplica
	PickupID    int64      `json:"pickup_id,omitempty"`   // ID del pickup si aplica
	DeliveryID  int64      `json:"delivery_id,omitempty"` // ID del delivery si aplica
}

// UnassignedPoint representa un punto no asignado
type UnassignedPoint struct {
	Location  [2]float64 `json:"location"`
	JobID     int64      `json:"job_id"`
	Reason    string     `json:"reason"`
	OrderRefs []string   `json:"order_refs,omitempty"` // ReferenceIDs de las órdenes no asignadas
}

// extractOrderRefsFromStep extrae las referenceID de las órdenes asociadas a un step
func extractOrderRefsFromStep(step Step, req optimization.FleetOptimization) []string {
	var orderRefs []string

	// Buscar órdenes asociadas al job
	if step.Job > 0 {
		visit := findVisitByJobID(step.Job, req.Visits)
		if visit != nil {
			for _, order := range visit.Orders {
				orderRefs = append(orderRefs, order.ReferenceID)
			}
		}
	}

	// Buscar órdenes asociadas al shipment
	if step.Shipment > 0 {
		visit := findVisitByShipmentID(step.Shipment, req.Visits)
		if visit != nil {
			for _, order := range visit.Orders {
				orderRefs = append(orderRefs, order.ReferenceID)
			}
		}
	}

	return orderRefs
}

// ExportToPolylineJSON exporta las rutas en formato optimizado para Leaflet
func (ret VroomOptimizationResponse) ExportToPolylineJSON(filename string, req optimization.FleetOptimization) error {
	var routesData []RouteData

	// Procesar cada ruta
	for i, route := range ret.Routes {
		routeData := RouteData{
			Vehicle:  route.Vehicle,
			Cost:     route.Cost,
			Duration: route.Duration,
		}

		// Decodificar polyline si existe
		if route.Geometry != "" {
			coords, _, err := polyline.DecodeCoords([]byte(route.Geometry))
			if err != nil {
				fmt.Printf("Error decodificando polyline para ruta %d: %v\n", i, err)
				continue
			}

			// Convertir a formato [lat, lng] para Leaflet
			for _, coord := range coords {
				routeData.Route = append(routeData.Route, []float64{coord[0], coord[1]}) // lat, lng
			}
		}

		// Procesar steps usando secuencia completa de la ruta
		stepCounter := 1 // Contador que empieza en 1 para jobs/pickups/deliveries
		for _, step := range route.Steps {
			if len(step.Location) == 2 {
				var stepNumber int
				if step.Type == "job" || step.Type == "pickup" || step.Type == "delivery" {
					stepNumber = stepCounter
					stepCounter++
				} else {
					stepNumber = 0 // Para 'start', 'end', 'break', etc.
				}
				stepPoint := StepPoint{
					Location:    [2]float64{step.Location[1], step.Location[0]}, // lat, lng
					StepType:    step.Type,
					StepNumber:  stepNumber, // Solo numerar jobs/pickups/deliveries
					Arrival:     step.Arrival,
					Description: step.Description,
					OrderRefs:   extractOrderRefsFromStep(step, req),
					JobID:       step.Job,
					ShipmentID:  step.Shipment,
					PickupID:    step.Pickup,
					DeliveryID:  step.Delivery,
				}
				routeData.Steps = append(routeData.Steps, stepPoint)
			}
		}

		routesData = append(routesData, routeData)
	}

	// Procesar trabajos no asignados
	for _, unassigned := range ret.Unassigned {
		if len(unassigned.Location) == 2 {
			// Buscar órdenes asociadas al trabajo no asignado
			var orderRefs []string
			visit := findVisitByJobID(unassigned.ID, req.Visits)
			if visit != nil {
				for _, order := range visit.Orders {
					orderRefs = append(orderRefs, order.ReferenceID)
				}
			} else {
				// Si no se encuentra como job, intentar como shipment
				visit = findVisitByShipmentID(unassigned.ID, req.Visits)
				if visit != nil {
					for _, order := range visit.Orders {
						orderRefs = append(orderRefs, order.ReferenceID)
					}
				}
			}

			unassignedPoint := UnassignedPoint{
				Location:  [2]float64{unassigned.Location[1], unassigned.Location[0]}, // lat, lng
				JobID:     unassigned.ID,
				Reason:    unassigned.Reason,
				OrderRefs: orderRefs,
			}
			if len(routesData) > 0 {
				routesData[0].Unassigned = append(routesData[0].Unassigned, unassignedPoint)
			}
		}
	}

	// Convertir a JSON
	jsonData, err := json.MarshalIndent(routesData, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando datos de ruta: %w", err)
	}

	// Asegurarse de que el directorio de destino exista
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creando directorio %s: %w", dir, err)
	}

	// Escribir al archivo
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	fmt.Printf("Datos de ruta exportados exitosamente a: %s\n", filename)
	fmt.Printf("Total de rutas: %d\n", len(routesData))

	return nil
}
