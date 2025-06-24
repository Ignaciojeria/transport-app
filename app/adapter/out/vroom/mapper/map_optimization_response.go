package mapper

import (
	"context"
	"fmt"
	"time"
	"transport-app/app/adapter/out/vroom/model"
	"transport-app/app/domain/optimization"
)

// MapOptimizationResponse convierte la respuesta de VROOM al modelo OptimizedFleet del dominio
func MapOptimizationResponse(
	ctx context.Context,
	vroomResponse model.VroomOptimizationResponse,
	originalFleet optimization.FleetOptimization,
) (optimization.OptimizedFleet, error) {
	// Crear mapeos para preservar la semántica de las visitas originales
	visitMappings := model.CreateVisitMappings(ctx, originalFleet.Visits)

	if err := vroomResponse.ExportToPolylineJSON("ui/static/dev/polyline.json"); err != nil {
		fmt.Printf("error exportando datos de ruta: %v\n", err)
	}
	// Mapear rutas optimizadas
	optimizedRoutes := make([]optimization.OptimizedRoute, 0, len(vroomResponse.Routes))
	for _, vroomRoute := range vroomResponse.Routes {
		// Encontrar el vehículo correspondiente
		var vehiclePlate string
		if int(vroomRoute.Vehicle) <= len(originalFleet.Vehicles) {
			vehiclePlate = originalFleet.Vehicles[vroomRoute.Vehicle-1].Plate
		} else {
			vehiclePlate = fmt.Sprintf("Vehicle-%d", vroomRoute.Vehicle)
		}

		// Mapear los pasos de la ruta
		optimizedSteps := make([]optimization.OptimizedStep, 0, len(vroomRoute.Steps))
		for _, vroomStep := range vroomRoute.Steps {
			optimizedStep := mapVroomStepToOptimizedStep(ctx, vroomStep, &visitMappings, originalFleet)
			if optimizedStep != nil {
				optimizedSteps = append(optimizedSteps, *optimizedStep)
			}
		}

		optimizedRoute := optimization.OptimizedRoute{
			VehiclePlate: vehiclePlate,
			Steps:        optimizedSteps,
			Cost:         vroomRoute.Cost,
			Duration:     vroomRoute.Duration,
		}

		optimizedRoutes = append(optimizedRoutes, optimizedRoute)
	}

	// Mapear elementos no asignados
	optimizedUnassigned := mapUnassignedToOptimizedUnassigned(
		ctx,
		vroomResponse.Unassigned,
		&visitMappings,
		originalFleet,
	)

	return optimization.OptimizedFleet{
		PlannedDate: time.Now(), // O usar una fecha específica si se proporciona
		Routes:      optimizedRoutes,
		Unassigned:  optimizedUnassigned,
	}, nil
}

// mapVroomStepToOptimizedStep convierte un paso de VROOM a un paso optimizado del dominio
func mapVroomStepToOptimizedStep(
	ctx context.Context,
	vroomStep model.Step,
	visitMappings *model.VisitMappings,
	originalFleet optimization.FleetOptimization,
) *optimization.OptimizedStep {
	// Determinar el tipo de paso
	stepType := mapVroomStepType(vroomStep.Type)

	var originalVisit *optimization.Visit
	visitIndex := -1 // Por defecto, -1 si no se encuentra

	if vroomStep.Job != 0 { // Es un Job (solo delivery, como en tu caso)
		jobIDToVisitIndex := visitMappings.GetJobIDToVisit()
		if index, exists := jobIDToVisitIndex[vroomStep.Job]; exists {
			if index < len(originalFleet.Visits) {
				originalVisit = &originalFleet.Visits[index]
				visitIndex = index
			}
		}
	} else if vroomStep.Shipment != 0 { // Es un Shipment (pickup y delivery)
		shipmentIDToVisit := visitMappings.GetShipmentIDToVisit()
		if visit, exists := shipmentIDToVisit[vroomStep.Shipment]; exists {
			originalVisit = visit
			// Buscar el índice del shipment de forma correcta
			for i := range originalFleet.Visits {
				if &originalFleet.Visits[i] == originalVisit {
					visitIndex = i
					break
				}
			}
		}
	}

	// Si no encontramos la visita, crear un paso básico
	if originalVisit == nil {
		return &optimization.OptimizedStep{
			Type:       stepType,
			VisitIndex: -1, // Indicar que no corresponde a una visita original
			Location: optimization.Location{
				Latitude:  vroomStep.Location[1], // VROOM usa [lon, lat], convertimos a [lat, lon]
				Longitude: vroomStep.Location[0],
				NodeInfo:  optimization.NodeInfo{},
			},
			Orders: []optimization.Order{},
		}
	}

	// Mapear ubicación basada en el tipo de paso
	var location optimization.Location
	switch stepType {
	case "pickup":
		location = optimization.Location{
			Latitude:  originalVisit.Pickup.Coordinates.Latitude,
			Longitude: originalVisit.Pickup.Coordinates.Longitude,
			NodeInfo:  originalVisit.Pickup.NodeInfo,
		}
	case "delivery":
		location = optimization.Location{
			Latitude:  originalVisit.Delivery.Coordinates.Latitude,
			Longitude: originalVisit.Delivery.Coordinates.Longitude,
			NodeInfo:  originalVisit.Delivery.NodeInfo,
		}
	default:
		// Para pasos start/end, usar la ubicación de VROOM si está disponible
		if len(vroomStep.Location) == 2 {
			location = optimization.Location{
				Latitude:  vroomStep.Location[1],
				Longitude: vroomStep.Location[0],
				NodeInfo:  optimization.NodeInfo{},
			}
		}
	}

	return &optimization.OptimizedStep{
		Type:       stepType,
		VisitIndex: visitIndex,
		Location:   location,
		Orders:     originalVisit.Orders,
	}
}

// mapVroomStepType convierte el tipo de paso de VROOM al tipo del dominio
func mapVroomStepType(vroomStepType string) string {
	switch vroomStepType {
	case "start":
		return "start"
	case "job":
		return "delivery"
	case "pickup":
		return "pickup"
	case "delivery":
		return "delivery"
	case "end":
		return "end"
	default:
		return vroomStepType
	}
}

// mapUnassignedToOptimizedUnassigned mapea los trabajos no asignados
func mapUnassignedToOptimizedUnassigned(
	ctx context.Context,
	unassignedJobs []model.UnassignedJob,
	visitMappings *model.VisitMappings,
	originalFleet optimization.FleetOptimization,
) optimization.OptimizedUnassigned {
	var unassignedOrders []optimization.Order
	var unassignedVehicles []optimization.Vehicle
	var unassignedOrigins []optimization.NodeInfo

	// Procesar trabajos no asignados
	for _, unassignedJob := range unassignedJobs {
		var originalVisit *optimization.Visit

		// Buscar en mapeos de jobs
		jobIDToVisitIndex := visitMappings.GetJobIDToVisit()
		if index, exists := jobIDToVisitIndex[unassignedJob.ID]; exists {
			if index < len(originalFleet.Visits) {
				originalVisit = &originalFleet.Visits[index]
			}
		}

		// Buscar en mapeos de shipments (si no se encontró como job)
		if originalVisit == nil {
			shipmentIDToVisit := visitMappings.GetShipmentIDToVisit()
			if visit, exists := shipmentIDToVisit[unassignedJob.ID]; exists {
				originalVisit = visit
			}
		}

		if originalVisit != nil {
			// Agregar las órdenes de la visita no asignada
			unassignedOrders = append(unassignedOrders, originalVisit.Orders...)

			// Agregar el origen si es un shipment
			if originalVisit.Pickup.Coordinates.Longitude != 0 || originalVisit.Pickup.Coordinates.Latitude != 0 {
				unassignedOrigins = append(unassignedOrigins, originalVisit.Pickup.NodeInfo)
			}
		}
	}

	return optimization.OptimizedUnassigned{
		Orders:   unassignedOrders,
		Vehicles: unassignedVehicles,
		Origins:  unassignedOrigins,
	}
}
