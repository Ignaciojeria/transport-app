package vroom

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"transport-app/app/domain/optimization"
)

// PolylineRouteData es una estructura desacoplada para exportar rutas a Leaflet o similar.
type PolylineRouteData struct {
	Route        [][]float64 // Puntos [lat, lng]
	VehiclePlate string
	Cost         int64
	Duration     int64
	Steps        []PolylineStepPoint
}

type PolylineStepPoint struct {
	Location     [2]float64
	StepType     string
	StepNumber   int
	ReferenceIDs []string
	IsEndPoint   bool
}

// GeneratePolylineDataFromOptimizedFleet genera la estructura de rutas y puntos igual a ExportToPolylineJSON, pero desde el dominio.
func GeneratePolylineDataFromOptimizedFleet(fleet optimization.OptimizedFleet) []PolylineRouteData {
	var routesData []PolylineRouteData

	for _, route := range fleet.Routes {
		routeData := PolylineRouteData{
			VehiclePlate: route.VehiclePlate,
			Cost:         route.Cost,
			Duration:     route.Duration,
		}

		var stepPoints []PolylineStepPoint
		var polyline [][]float64

		for j, step := range route.Steps {
			lat := step.Location.Latitude
			lng := step.Location.Longitude
			location := [2]float64{lat, lng}
			if lat != 0 || lng != 0 {
				polyline = append(polyline, []float64{lat, lng})
			}

			// Puedes agregar lógica para ReferenceIDs si tu dominio los tiene
			var referenceIDs []string
			for _, order := range step.Orders {
				referenceIDs = append(referenceIDs, order.ReferenceID)
			}

			isEndPoint := j == len(route.Steps)-1

			stepPoints = append(stepPoints, PolylineStepPoint{
				Location:     location,
				StepType:     step.Type,
				StepNumber:   j + 1,
				ReferenceIDs: referenceIDs,
				IsEndPoint:   isEndPoint,
			})
		}

		routeData.Route = polyline
		routeData.Steps = stepPoints
		routesData = append(routesData, routeData)
	}

	return routesData
}

// ExportPolylineDataToFile exporta el polyline y los steps a un archivo JSON.
func ExportPolylineDataToFile(filename string, routesData []PolylineRouteData) error {
	jsonData, err := json.MarshalIndent(routesData, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando datos de polyline: %w", err)
	}

	dir := ""
	if idx := len(filename) - len("/"); idx > 0 {
		dir = filename[:idx]
	}
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creando directorio %s: %w", dir, err)
		}
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	fmt.Printf("Polyline exportado exitosamente a: %s\n", filename)
	return nil
}

type RouteData struct {
	Route        [][]float64       `json:"route"`
	Steps        []StepPoint       `json:"steps"`
	VehiclePlate string            `json:"vehicle_plate"`
	Cost         int64             `json:"cost"`
	Duration     int64             `json:"duration"`
	Unassigned   []UnassignedPoint `json:"unassigned,omitempty"`
}

type StepPoint struct {
	Location             [2]float64 `json:"location"`
	StepType             string     `json:"step_type"`
	StepNumber           int        `json:"step_number"`
	Arrival              int64      `json:"arrival"`
	EstimatedArrivalTime string     `json:"estimated_arrival_time"`
	Description          string     `json:"description,omitempty"`
	ReferenceIDs         []string   `json:"reference_ids,omitempty"`
	IsEndPoint           bool       `json:"is_end_point,omitempty"`
	VehiclePlate         string     `json:"vehicle_plate,omitempty"`
	TotalCost            int64      `json:"total_cost,omitempty"`
	TotalDuration        int64      `json:"total_duration,omitempty"`
}

type UnassignedPoint struct {
	Location [2]float64 `json:"location"`
	JobID    string     `json:"job_id"`
	Reason   string     `json:"reason"`
}

// ExportPolylineJSONFromOptimizedFleet exporta rutas y steps desde OptimizedFleet en formato idéntico a ExportToPolylineJSON.
// Usa la hora de inicio específica de cada vehículo desde el FleetOptimization original.
func ExportPolylineJSONFromOptimizedFleet(filename string, fleet optimization.OptimizedFleet, originalFleet optimization.FleetOptimization) error {
	// Crear mapa de patentes a horas de inicio
	vehicleStartHours := make(map[string]string)
	for _, v := range originalFleet.Vehicles {
		vehicleStartHours[v.Plate] = v.TimeWindow.Start
	}

	return ExportPolylineJSONFromOptimizedFleetWithStartHours(filename, fleet, vehicleStartHours)
}

// ExportPolylineJSONFromOptimizedFleetWithStartHours exporta rutas y steps desde OptimizedFleet usando un mapa de patentes a horas de inicio.
// vehicleStartHours es un mapa donde la clave es la patente del vehículo y el valor es la hora de inicio en formato "HH:MM"
func ExportPolylineJSONFromOptimizedFleetWithStartHours(filename string, fleet optimization.OptimizedFleet, vehicleStartHours map[string]string) error {
	var routesData []RouteData

	for _, route := range fleet.Routes {
		// Obtener la hora de inicio del vehículo desde el mapa
		startHour, exists := vehicleStartHours[route.VehiclePlate]
		if !exists || startHour == "" {
			startHour = "08:00" // Hora por defecto si no está definida
		}

		startTime, err := time.Parse("15:04", startHour)
		if err != nil {
			return fmt.Errorf("hora de inicio inválida para vehículo %s: %w", route.VehiclePlate, err)
		}

		routeData := RouteData{
			VehiclePlate: route.VehiclePlate,
			Cost:         route.Cost,
			Duration:     route.Duration,
		}

		var polyline [][]float64
		var steps []StepPoint
		for j, step := range route.Steps {
			lat := step.Location.Latitude
			lng := step.Location.Longitude
			location := [2]float64{lat, lng}

			var referenceIDs []string
			for _, order := range step.Orders {
				referenceIDs = append(referenceIDs, order.ReferenceID)
			}

			isEndPoint := j == len(route.Steps)-1

			// Lógica especial para el último end
			if isEndPoint && step.Type == "end" && j > 0 {
				prevStep := route.Steps[j-1]
				if prevStep.Location.Latitude == step.Location.Latitude && prevStep.Location.Longitude == step.Location.Longitude &&
					(prevStep.Type == "delivery" || prevStep.Type == "job") {
					continue
				}
			}

			// Incluir siempre pickups, deliveries, start y end
			if step.Type == "pickup" || step.Type == "delivery" || step.Type == "start" || step.Type == "end" {
				polyline = append(polyline, []float64{lat, lng})

				// Usar SequenceNumber de la primera orden si está disponible
				stepNumber := 0
				if len(step.Orders) > 0 && step.Orders[0].SequenceNumber != nil {
					stepNumber = *step.Orders[0].SequenceNumber
				}

				estimatedArrival := startTime.Add(time.Duration(step.Arrival) * time.Second).Format("15:04")

				steps = append(steps, StepPoint{
					Location:             location,
					StepType:             step.Type,
					StepNumber:           stepNumber,
					Arrival:              step.Arrival,
					EstimatedArrivalTime: estimatedArrival,
					ReferenceIDs:         referenceIDs,
					IsEndPoint:           isEndPoint,
					VehiclePlate:         route.VehiclePlate,
					TotalCost:            route.Cost,
					TotalDuration:        route.Duration,
				})
			}
		}
		routeData.Route = polyline
		routeData.Steps = steps
		routesData = append(routesData, routeData)
	}

	jsonData, err := json.MarshalIndent(routesData, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando datos de polyline: %w", err)
	}

	dir := ""
	if idx := len(filename) - len("/"); idx > 0 {
		dir = filename[:idx]
	}
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creando directorio %s: %w", dir, err)
		}
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	fmt.Printf("Polyline exportado exitosamente a: %s\n", filename)
	return nil
}
