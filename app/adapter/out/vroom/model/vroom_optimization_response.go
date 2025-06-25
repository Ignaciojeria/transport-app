package model

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"transport-app/app/domain"
	"transport-app/app/domain/optimization"

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

// VisitMappings contiene los mapeos bidireccionales para preservar la semántica
type VisitMappings struct {
	// Mapeo de Job ID de VROOM a índice de la Visit original
	jobIDToVisit map[int64]int
	// Mapeo de Shipment ID de VROOM a Visit original
	shipmentIDToVisit map[int64]*optimization.Visit
	// Mapeo inverso para debugging
	visitToJobID      map[string]int64
	visitToShipmentID map[string]int64
}

// GetJobIDToVisit retorna el mapeo de Job ID a índice de Visit
func (vm *VisitMappings) GetJobIDToVisit() map[int64]int {
	return vm.jobIDToVisit
}

// GetShipmentIDToVisit retorna el mapeo de Shipment ID a Visit
func (vm *VisitMappings) GetShipmentIDToVisit() map[int64]*optimization.Visit {
	return vm.shipmentIDToVisit
}

// CreateVisitMappings crea los mapeos bidireccionales basados en la lógica de mapeo de VROOM
func CreateVisitMappings(ctx context.Context, visits []optimization.Visit) VisitMappings {
	mappings := VisitMappings{
		jobIDToVisit:      make(map[int64]int),
		shipmentIDToVisit: make(map[int64]*optimization.Visit),
		visitToJobID:      make(map[string]int64),
		visitToShipmentID: make(map[string]int64),
	}

	// Crear registros que simulen exactamente la lógica del mapper de request
	locationRegistry := newLocationContactRegistry()

	for i, visit := range visits {
		// Verificar si tenemos pickup válido
		hasValidPickup := visit.Pickup.Coordinates.Longitude != 0 || visit.Pickup.Coordinates.Latitude != 0
		// Verificar si tenemos delivery válido
		hasValidDelivery := visit.Delivery.Coordinates.Longitude != 0 || visit.Delivery.Coordinates.Latitude != 0

		// Crear identificador único para la visita basado en sus características
		visitKey := createVisitKey(visit)

		if hasValidPickup && hasValidDelivery {
			// Ambos son válidos -> es un Shipment
			shipmentID := int64(i + 1) // VROOM usa índices basados en 1
			mappings.shipmentIDToVisit[shipmentID] = &visits[i]
			mappings.visitToShipmentID[visitKey] = shipmentID
		} else if hasValidDelivery {
			// Solo delivery válido -> es un Job
			// Usar exactamente la misma lógica que en el mapper de request
			deliveryLocationKey := generateLocationKey(visit.Delivery.Coordinates.Latitude, visit.Delivery.Coordinates.Longitude)
			deliveryContactID := getContactID(ctx, visit.Delivery.Contact)
			jobID := locationRegistry.getLocationContactID(deliveryLocationKey, deliveryContactID)

			mappings.jobIDToVisit[jobID] = i
			mappings.visitToJobID[visitKey] = jobID
		}
	}

	return mappings
}

// locationContactRegistry simula exactamente la lógica del mapper de request
type locationContactRegistry struct {
	counter int64
	mapping map[string]int64
}

func newLocationContactRegistry() *locationContactRegistry {
	return &locationContactRegistry{
		counter: 1,
		mapping: make(map[string]int64),
	}
}

func (r *locationContactRegistry) getLocationContactID(coordinates string, contactID string) int64 {
	key := fmt.Sprintf("%s_%s", coordinates, contactID)
	if id, ok := r.mapping[key]; ok {
		return id
	}
	r.mapping[key] = r.counter
	r.counter++
	return r.mapping[key]
}

// generateLocationKey debe usar exactamente el mismo formato que en el mapper de request
func generateLocationKey(lat, lon float64) string {
	return fmt.Sprintf("%.6f_%.6f", lat, lon)
}

// getContactID debe usar exactamente la misma lógica que en el mapper de request
func getContactID(ctx context.Context, contact optimization.Contact) string {
	// Si no hay información de contacto, usar coordenadas como identificador único
	if contact.Email == "" && contact.Phone == "" && contact.NationalID == "" && contact.FullName == "" {
		return ""
	}

	// Usar la función DocID del dominio Contact
	domainContact := domain.Contact{
		PrimaryEmail: contact.Email,
		PrimaryPhone: contact.Phone,
		NationalID:   contact.NationalID,
		FullName:     contact.FullName,
	}

	docID := domainContact.DocID(ctx)
	return string(docID)
}

// createVisitKey crea una clave única para identificar una visita
func createVisitKey(visit optimization.Visit) string {
	// Usar las coordenadas de delivery como identificador principal
	deliveryKey := fmt.Sprintf("%.6f,%.6f",
		visit.Delivery.Coordinates.Latitude,
		visit.Delivery.Coordinates.Longitude)

	// Si hay pickup válido, incluirlo en la clave
	if visit.Pickup.Coordinates.Longitude != 0 || visit.Pickup.Coordinates.Latitude != 0 {
		pickupKey := fmt.Sprintf("%.6f,%.6f",
			visit.Pickup.Coordinates.Latitude,
			visit.Pickup.Coordinates.Longitude)
		return fmt.Sprintf("shipment:%s->%s", pickupKey, deliveryKey)
	}

	return fmt.Sprintf("job:%s", deliveryKey)
}

// createOrdersFromVisit crea órdenes del dominio basadas en una visita
func createOrdersFromVisit(visit *optimization.Visit, hasPickup bool) []domain.Order {
	var orders []domain.Order

	// Crear orden para cada order en la visita
	for _, orderReq := range visit.Orders {
		order := domain.Order{
			ReferenceID:    domain.ReferenceID(orderReq.ReferenceID),
			SequenceNumber: orderReq.SequenceNumber,
			Destination: domain.NodeInfo{
				ReferenceID: domain.ReferenceID(visit.Delivery.NodeInfo.ReferenceID),
				Name:        "Destino de entrega",
				AddressInfo: domain.AddressInfo{
					Coordinates: domain.Coordinates{
						Point: orb.Point{visit.Delivery.Coordinates.Longitude, visit.Delivery.Coordinates.Latitude},
					},
					Contact: domain.Contact{
						FullName:     visit.Delivery.Contact.FullName,
						PrimaryEmail: visit.Delivery.Contact.Email,
						PrimaryPhone: visit.Delivery.Contact.Phone,
						NationalID:   visit.Delivery.Contact.NationalID,
					},
				},
			},
		}

		// Para shipments (pickup + delivery), incluir origen
		if hasPickup {
			order.Origin = domain.NodeInfo{
				ReferenceID: domain.ReferenceID(visit.Pickup.NodeInfo.ReferenceID),
				Name:        "Origen de recogida",
				AddressInfo: domain.AddressInfo{
					Coordinates: domain.Coordinates{
						Point: orb.Point{visit.Pickup.Coordinates.Longitude, visit.Pickup.Coordinates.Latitude},
					},
					Contact: domain.Contact{
						FullName:     visit.Pickup.Contact.FullName,
						PrimaryEmail: visit.Pickup.Contact.Email,
						PrimaryPhone: visit.Pickup.Contact.Phone,
						NationalID:   visit.Pickup.Contact.NationalID,
					},
				},
			}
		}

		// Mapear delivery units
		var deliveryUnits domain.DeliveryUnits
		for _, duReq := range orderReq.DeliveryUnits {
			deliveryUnit := domain.DeliveryUnit{
				Lpn: duReq.Lpn,
				Weight: domain.Weight{
					Value: duReq.Weight,
					Unit:  "g",
				},
				Insurance: domain.Insurance{
					UnitValue: duReq.Insurance,
					Currency:  "CLP",
				},
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

// RouteData representa los datos de ruta para el frontend
type RouteData struct {
	Route        [][]float64       `json:"route"` // Coordenadas decodificadas del polyline
	Steps        []StepPoint       `json:"steps"` // Puntos de parada
	Vehicle      int64             `json:"vehicle"`
	VehiclePlate string            `json:"vehicle_plate"` // Patente del vehículo
	Cost         int64             `json:"cost"`
	Duration     int64             `json:"duration"`
	Unassigned   []UnassignedPoint `json:"unassigned,omitempty"`
}

// StepPoint representa un punto de parada
type StepPoint struct {
	Location     [2]float64 `json:"location"`
	StepType     string     `json:"step_type"`
	StepNumber   int        `json:"step_number"`
	Arrival      int64      `json:"arrival"`
	Description  string     `json:"description,omitempty"`
	ReferenceIDs []string   `json:"reference_ids,omitempty"` // ReferenceIDs de las órdenes asociadas
	// Metadata adicional para el punto final
	IsEndPoint    bool   `json:"is_end_point,omitempty"`   // Indica si es el punto final de la ruta
	VehicleID     int64  `json:"vehicle_id,omitempty"`     // ID del vehículo
	VehiclePlate  string `json:"vehicle_plate,omitempty"`  // Patente del vehículo
	TotalCost     int64  `json:"total_cost,omitempty"`     // Costo total de la ruta
	TotalDuration int64  `json:"total_duration,omitempty"` // Duración total de la ruta
}

// UnassignedPoint representa un punto no asignado
type UnassignedPoint struct {
	Location [2]float64 `json:"location"`
	JobID    int64      `json:"job_id"`
	Reason   string     `json:"reason"`
}

// ExportToPolylineJSON exporta las rutas en formato optimizado para Leaflet
func (ret VroomOptimizationResponse) ExportToPolylineJSON(filename string, originalFleet *optimization.FleetOptimization) error {
	var routesData []RouteData

	// Crear mapeos para preservar la semántica de las visitas originales
	var visitMappings VisitMappings
	if originalFleet != nil {
		visitMappings = CreateVisitMappings(context.Background(), originalFleet.Visits)
	}

	// Procesar cada ruta
	for i, route := range ret.Routes {
		// Obtener la patente del vehículo
		var vehiclePlate string
		if originalFleet != nil && int(route.Vehicle) <= len(originalFleet.Vehicles) {
			vehiclePlate = originalFleet.Vehicles[route.Vehicle-1].Plate
		} else {
			vehiclePlate = fmt.Sprintf("Vehicle-%d", route.Vehicle)
		}

		routeData := RouteData{
			Vehicle:      route.Vehicle,
			VehiclePlate: vehiclePlate,
			Cost:         route.Cost,
			Duration:     route.Duration,
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

		// Procesar steps
		pickupCount := 1
		deliveryCount := 1
		jobCount := 1

		for j, step := range route.Steps {
			// Determinar número secuencial por tipo
			stepNumber := j
			switch step.Type {
			case "pickup":
				stepNumber = pickupCount
				pickupCount++
			case "delivery":
				stepNumber = deliveryCount
				deliveryCount++
			case "job":
				stepNumber = jobCount
				jobCount++
			}

			// Obtener ReferenceIDs de las órdenes asociadas
			var referenceIDs []string
			if originalFleet != nil {
				referenceIDs = getOrderReferenceIDs(step, &visitMappings, originalFleet)
			}

			// Si es el último step y es de tipo end, solo incluirlo si NO hay un delivery/job en la misma ubicación justo antes
			isEndPoint := j == len(route.Steps)-1 // Último step de la ruta
			if isEndPoint && step.Type == "end" && j > 0 {
				prevStep := route.Steps[j-1]
				// Comparar ubicación con el paso anterior
				if len(step.Location) == 2 && len(prevStep.Location) == 2 &&
					step.Location[0] == prevStep.Location[0] && step.Location[1] == prevStep.Location[1] &&
					(prevStep.Type == "delivery" || prevStep.Type == "job") {
					// Si el paso anterior es una entrega/job en la misma ubicación, NO incluir el end
					continue
				}
			}

			// Solo incluir steps que tengan órdenes asociadas o sean start/end
			if len(referenceIDs) > 0 || step.Type == "start" || step.Type == "end" {
				var location [2]float64

				// Usar coordenadas de VROOM si están disponibles
				if len(step.Location) == 2 {
					location = [2]float64{step.Location[1], step.Location[0]} // lat, lng
				} else {
					// Para steps sin coordenadas directas, intentar obtenerlas de la visita original
					if originalFleet != nil {
						var originalVisit *optimization.Visit

						if step.Job != 0 { // Es un Job (solo delivery)
							jobIDToVisitIndex := visitMappings.GetJobIDToVisit()
							if index, exists := jobIDToVisitIndex[step.Job]; exists {
								if index < len(originalFleet.Visits) {
									originalVisit = &originalFleet.Visits[index]
								}
							}
						} else if step.Shipment != 0 { // Es un Shipment (pickup y delivery)
							shipmentIDToVisit := visitMappings.GetShipmentIDToVisit()
							if visit, exists := shipmentIDToVisit[step.Shipment]; exists {
								originalVisit = visit
							}
						}

						if originalVisit != nil {
							switch step.Type {
							case "pickup":
								location = [2]float64{originalVisit.Pickup.Coordinates.Latitude, originalVisit.Pickup.Coordinates.Longitude}
							case "delivery", "job":
								location = [2]float64{originalVisit.Delivery.Coordinates.Latitude, originalVisit.Delivery.Coordinates.Longitude}
							}
						}
					}
				}

				stepPoint := StepPoint{
					Location:     location,
					StepType:     step.Type,
					StepNumber:   stepNumber,
					Arrival:      step.Arrival,
					Description:  step.Description,
					ReferenceIDs: referenceIDs,
					// Metadata adicional para el punto final
					IsEndPoint:    isEndPoint,
					VehicleID:     route.Vehicle,
					VehiclePlate:  vehiclePlate,
					TotalCost:     route.Cost,
					TotalDuration: route.Duration,
				}
				routeData.Steps = append(routeData.Steps, stepPoint)
			}
		}

		routesData = append(routesData, routeData)
	}

	// Procesar trabajos no asignados
	for _, unassigned := range ret.Unassigned {
		if len(unassigned.Location) == 2 {
			unassignedPoint := UnassignedPoint{
				Location: [2]float64{unassigned.Location[1], unassigned.Location[0]}, // lat, lng
				JobID:    unassigned.ID,
				Reason:   unassigned.Reason,
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

// getOrderReferenceIDs obtiene los ReferenceIDs de las órdenes asociadas a un step
func getOrderReferenceIDs(step Step, visitMappings *VisitMappings, originalFleet *optimization.FleetOptimization) []string {
	var referenceIDs []string

	var originalVisit *optimization.Visit
	if step.Job != 0 { // Es un Job (solo delivery)
		jobIDToVisitIndex := visitMappings.GetJobIDToVisit()
		if index, exists := jobIDToVisitIndex[step.Job]; exists {
			if index < len(originalFleet.Visits) {
				originalVisit = &originalFleet.Visits[index]
			}
		}
	} else if step.Shipment != 0 { // Es un Shipment (pickup y delivery)
		shipmentIDToVisit := visitMappings.GetShipmentIDToVisit()
		if visit, exists := shipmentIDToVisit[step.Shipment]; exists {
			originalVisit = visit
		}
	}

	if originalVisit != nil {
		for _, order := range originalVisit.Orders {
			referenceIDs = append(referenceIDs, string(order.ReferenceID))
		}
	}

	return referenceIDs
}
