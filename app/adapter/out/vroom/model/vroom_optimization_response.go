package model

import (
	"context"
	"fmt"
	"transport-app/app/domain/optimization"
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

	// Usar exactamente la misma lógica secuencial que el mapper de request
	jobCounter := int64(1)

	for i, visit := range visits {
		// Verificar si tenemos pickup válido usando la misma lógica que el mapper de request
		hasValidPickup := isValidCoordinates(visit.Pickup.Coordinates.Latitude, visit.Pickup.Coordinates.Longitude)
		// Verificar si tenemos delivery válido usando la misma lógica que el mapper de request
		hasValidDelivery := isValidCoordinates(visit.Delivery.Coordinates.Latitude, visit.Delivery.Coordinates.Longitude)

		// Crear identificador único para la visita basado en sus características
		visitKey := createVisitKey(visit)

		if hasValidPickup && hasValidDelivery {
			// Ambos son válidos -> es un Shipment
			// Para shipments, VROOM usa el índice de la visita + 1 como ID del shipment
			shipmentID := int64(i + 1)
			mappings.shipmentIDToVisit[shipmentID] = &visits[i]
			mappings.visitToShipmentID[visitKey] = shipmentID
			fmt.Printf("DEBUG: Creando Shipment ID %d para visita %d\n", shipmentID, i)

			// Los pasos de pickup y delivery de un shipment usan jobCounter secuencial
			// pickup step ID = jobCounter
			// delivery step ID = jobCounter + 1
			pickupJobID := jobCounter
			deliveryJobID := jobCounter + 1
			jobCounter += 2

			// Mapear ambos pasos a la misma visita
			mappings.jobIDToVisit[pickupJobID] = i
			mappings.jobIDToVisit[deliveryJobID] = i
			fmt.Printf("DEBUG: Mapeando Job IDs %d (pickup) y %d (delivery) para Shipment ID %d\n", pickupJobID, deliveryJobID, shipmentID)

		} else if hasValidPickup {
			// Solo pickup válido -> es un Job
			jobID := jobCounter
			jobCounter++

			mappings.jobIDToVisit[jobID] = i
			mappings.visitToJobID[visitKey] = jobID
			fmt.Printf("DEBUG: Creando Job ID %d para visita %d (pickup)\n", jobID, i)

		} else if hasValidDelivery {
			// Solo delivery válido -> es un Job
			jobID := jobCounter
			jobCounter++

			mappings.jobIDToVisit[jobID] = i
			mappings.visitToJobID[visitKey] = jobID
			fmt.Printf("DEBUG: Creando Job ID %d para visita %d (delivery)\n", jobID, i)
		}
	}

	return mappings
}

// createVisitKey crea una clave única para identificar una visita
func createVisitKey(visit optimization.Visit) string {
	// Usar las coordenadas de delivery como identificador principal
	deliveryKey := fmt.Sprintf("%.6f,%.6f",
		visit.Delivery.Coordinates.Latitude,
		visit.Delivery.Coordinates.Longitude)

	// Si hay pickup válido, incluirlo en la clave
	if isValidCoordinates(visit.Pickup.Coordinates.Latitude, visit.Pickup.Coordinates.Longitude) {
		pickupKey := fmt.Sprintf("%.6f,%.6f",
			visit.Pickup.Coordinates.Latitude,
			visit.Pickup.Coordinates.Longitude)
		return fmt.Sprintf("shipment:%s->%s", pickupKey, deliveryKey)
	}

	return fmt.Sprintf("job:%s", deliveryKey)
}

// isValidCoordinates valida que las coordenadas sean válidas
// Debe usar exactamente la misma lógica que en el mapper de request
func isValidCoordinates(lat, lon float64) bool {
	return lat != 0 && lon != 0
}
