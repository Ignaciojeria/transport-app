package model

import (
	"context"
	"fmt"
	"transport-app/app/domain"
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
