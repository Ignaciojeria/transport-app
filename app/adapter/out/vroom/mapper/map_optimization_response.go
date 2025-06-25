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
		sequence := 1

		for stepIndex, vroomStep := range vroomRoute.Steps {
			// Determinar la secuencia basada en el orden de los pasos con órdenes
			var stepSequence int
			if vroomStep.Type == "pickup" || vroomStep.Type == "delivery" || vroomStep.Type == "job" {
				// Cada paso con órdenes recibe un número secuencial
				stepSequence = sequence
				sequence++
			} else {
				// Para pasos start/end, no asignar secuencia
				stepSequence = 0
			}

			optimizedStep, err := mapVroomStepToOptimizedStep(ctx, vroomStep, &visitMappings, originalFleet, stepIndex, stepSequence)
			if err != nil {
				return optimization.OptimizedFleet{}, err
			}
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

	optimizedFleet := optimization.OptimizedFleet{
		PlannedDate: time.Now(), // O usar una fecha específica si se proporciona
		Routes:      optimizedRoutes,
		Unassigned:  optimizedUnassigned,
	}

	return optimizedFleet, nil
}

// mapVroomStepToOptimizedStep convierte un paso de VROOM a un paso optimizado del dominio
func mapVroomStepToOptimizedStep(
	ctx context.Context,
	vroomStep model.Step,
	visitMappings *model.VisitMappings,
	originalFleet optimization.FleetOptimization,
	stepIndex int,
	stepSequence int,
) (*optimization.OptimizedStep, error) {
	var originalVisit *optimization.Visit
	var visitIndex int = -1

	// Buscar la visita original basada en el tipo de paso
	if vroomStep.Job != 0 { // Es un Job (solo delivery, como en tu caso)
		// PRIMERO: Intentar buscar si este job corresponde a un shipment original
		shipmentIDToVisit := visitMappings.GetShipmentIDToVisit()

		// Buscar si el jobID corresponde directamente a un shipment
		if visit, exists := shipmentIDToVisit[vroomStep.Job]; exists {
			originalVisit = visit
			// Buscar el índice de la visita en el slice original por contenido
			for i, v := range originalFleet.Visits {
				if v.Pickup.Coordinates.Latitude == visit.Pickup.Coordinates.Latitude &&
					v.Pickup.Coordinates.Longitude == visit.Pickup.Coordinates.Longitude &&
					v.Delivery.Coordinates.Latitude == visit.Delivery.Coordinates.Latitude &&
					v.Delivery.Coordinates.Longitude == visit.Delivery.Coordinates.Longitude {
					visitIndex = i
					break
				}
			}
			fmt.Printf("DEBUG: Encontrado Shipment ID %d -> Visita %d (por jobID directo)\n", vroomStep.Job, visitIndex)
		} else {
			// Buscar si el jobID corresponde a un paso de pickup o delivery de un shipment
			// Los pasos de pickup y delivery de un shipment tienen IDs consecutivos
			for shipmentID, visit := range shipmentIDToVisit {
				// Verificar si el jobID es el pickup o delivery del shipment
				if vroomStep.Job == shipmentID || vroomStep.Job == shipmentID+1 {
					originalVisit = visit
					// Buscar el índice de la visita en el slice original por contenido
					for i, v := range originalFleet.Visits {
						if v.Pickup.Coordinates.Latitude == visit.Pickup.Coordinates.Latitude &&
							v.Pickup.Coordinates.Longitude == visit.Pickup.Coordinates.Longitude &&
							v.Delivery.Coordinates.Latitude == visit.Delivery.Coordinates.Latitude &&
							v.Delivery.Coordinates.Longitude == visit.Delivery.Coordinates.Longitude {
							visitIndex = i
							break
						}
					}
					fmt.Printf("DEBUG: Encontrado Shipment ID %d -> Visita %d (jobID %d es pickup/delivery)\n", shipmentID, visitIndex, vroomStep.Job)
					break
				}
			}
		}

		// SEGUNDO: Si no se encontró como shipment, buscar en el mapeo normal de jobs
		if originalVisit == nil {
			jobIDToVisitIndex := visitMappings.GetJobIDToVisit()
			if index, exists := jobIDToVisitIndex[vroomStep.Job]; exists {
				if index < len(originalFleet.Visits) {
					originalVisit = &originalFleet.Visits[index]
					visitIndex = index
					fmt.Printf("DEBUG: Encontrado Job ID %d -> Visita %d\n", vroomStep.Job, index)
				}
			} else {
				fmt.Printf("DEBUG: Job ID %d NO encontrado en jobIDToVisitIndex\n", vroomStep.Job)
			}
		}
	} else if vroomStep.Shipment != 0 { // Es un Shipment (pickup y delivery)
		shipmentIDToVisit := visitMappings.GetShipmentIDToVisit()
		if visit, exists := shipmentIDToVisit[vroomStep.Shipment]; exists {
			originalVisit = visit
			// Buscar el índice de la visita en el slice original
			for i, v := range originalFleet.Visits {
				if &v == visit {
					visitIndex = i
					break
				}
			}
			fmt.Printf("DEBUG: Encontrado Shipment ID %d -> Visita %d\n", vroomStep.Shipment, visitIndex)
		} else {
			fmt.Printf("DEBUG: Shipment ID %d NO encontrado en shipmentIDToVisit\n", vroomStep.Shipment)
		}
	}

	if originalVisit == nil {
		// Para pasos start/end, no necesitamos una visita original
		if vroomStep.Type == "start" || vroomStep.Type == "end" {
			fmt.Printf("DEBUG: Paso tipo %s sin visita original (normal para start/end)\n", vroomStep.Type)
			// Determinar el tipo de paso
			stepType := mapVroomStepType(vroomStep.Type)

			// Crear ubicación básica para pasos start/end
			var location optimization.Location
			if len(vroomStep.Location) == 2 {
				location = optimization.Location{
					Latitude:  vroomStep.Location[1],
					Longitude: vroomStep.Location[0],
					NodeInfo:  optimization.NodeInfo{},
				}
			}

			return &optimization.OptimizedStep{
				Type:       stepType,
				Location:   location,
				VisitIndex: -1, // Indica que no corresponde a una visita específica
			}, nil
		}

		// Intentar buscar si este job corresponde a un shipment original
		if vroomStep.Job != 0 {
			shipmentIDToVisit := visitMappings.GetShipmentIDToVisit()

			// Buscar si el jobID corresponde a un shipment
			if visit, exists := shipmentIDToVisit[vroomStep.Job]; exists {
				originalVisit = visit
				// Buscar el índice de la visita en el slice original por contenido
				for i, v := range originalFleet.Visits {
					if v.Pickup.Coordinates.Latitude == visit.Pickup.Coordinates.Latitude &&
						v.Pickup.Coordinates.Longitude == visit.Pickup.Coordinates.Longitude &&
						v.Delivery.Coordinates.Latitude == visit.Delivery.Coordinates.Latitude &&
						v.Delivery.Coordinates.Longitude == visit.Delivery.Coordinates.Longitude {
						visitIndex = i
						break
					}
				}
				fmt.Printf("DEBUG: Encontrado Shipment ID %d -> Visita %d (por jobID)\n", vroomStep.Job, visitIndex)
			} else {
				// Buscar si el jobID corresponde a un paso de pickup o delivery de un shipment
				// Los pasos de pickup y delivery de un shipment tienen IDs consecutivos
				for shipmentID, visit := range shipmentIDToVisit {
					// Verificar si el jobID es el pickup o delivery del shipment
					if vroomStep.Job == shipmentID || vroomStep.Job == shipmentID+1 {
						originalVisit = visit
						// Buscar el índice de la visita en el slice original por contenido
						for i, v := range originalFleet.Visits {
							if v.Pickup.Coordinates.Latitude == visit.Pickup.Coordinates.Latitude &&
								v.Pickup.Coordinates.Longitude == visit.Pickup.Coordinates.Longitude &&
								v.Delivery.Coordinates.Latitude == visit.Delivery.Coordinates.Latitude &&
								v.Delivery.Coordinates.Longitude == visit.Delivery.Coordinates.Longitude {
								visitIndex = i
								break
							}
						}
						fmt.Printf("DEBUG: Encontrado Shipment ID %d -> Visita %d (jobID %d es pickup/delivery)\n", shipmentID, visitIndex, vroomStep.Job)
						break
					}
				}
			}
		}

		// Si aún no se encontró la visita, crear un paso básico sin visita
		if originalVisit == nil {
			fmt.Printf("DEBUG: Creando paso sin visita original para tipo %s, Job: %d, Shipment: %d\n", vroomStep.Type, vroomStep.Job, vroomStep.Shipment)

			var location optimization.Location
			if len(vroomStep.Location) == 2 {
				location = optimization.Location{
					Latitude:  vroomStep.Location[1],
					Longitude: vroomStep.Location[0],
					NodeInfo:  optimization.NodeInfo{},
				}
			}

			return &optimization.OptimizedStep{
				Type:       mapVroomStepType(vroomStep.Type),
				Location:   location,
				VisitIndex: -1, // Indica que no corresponde a una visita específica
			}, nil
		}
	}

	// Determinar el tipo de paso
	stepType := mapVroomStepType(vroomStep.Type)

	// Mapear ubicación basada en el tipo de paso
	var location optimization.Location
	switch stepType {
	case "pickup":
		// Siempre usar las coordenadas de VROOM si están disponibles
		if len(vroomStep.Location) == 2 {
			location = optimization.Location{
				Latitude:  vroomStep.Location[1], // VROOM usa [lon, lat], convertimos a [lat, lon]
				Longitude: vroomStep.Location[0],
				NodeInfo:  originalVisit.Pickup.NodeInfo,
			}
			fmt.Printf("DEBUG: Pickup para visita %d - usando coordenadas de VROOM: lat=%f, lon=%f\n", visitIndex, location.Latitude, location.Longitude)
		} else {
			// Si VROOM no tiene coordenadas, usar las originales
			location = optimization.Location{
				Latitude:  originalVisit.Pickup.Coordinates.Latitude,
				Longitude: originalVisit.Pickup.Coordinates.Longitude,
				NodeInfo:  originalVisit.Pickup.NodeInfo,
			}
			fmt.Printf("DEBUG: Pickup para visita %d - VROOM sin coordenadas, usando originales: lat=%f, lon=%f\n", visitIndex, location.Latitude, location.Longitude)
		}
	case "delivery":
		// Siempre usar las coordenadas de VROOM si están disponibles
		if len(vroomStep.Location) == 2 {
			location = optimization.Location{
				Latitude:  vroomStep.Location[1], // VROOM usa [lon, lat], convertimos a [lat, lon]
				Longitude: vroomStep.Location[0],
				NodeInfo:  originalVisit.Delivery.NodeInfo,
			}
			fmt.Printf("DEBUG: Delivery para visita %d - usando coordenadas de VROOM: lat=%f, lon=%f\n", visitIndex, location.Latitude, location.Longitude)
		} else {
			// Si VROOM no tiene coordenadas, usar las originales
			location = optimization.Location{
				Latitude:  originalVisit.Delivery.Coordinates.Latitude,
				Longitude: originalVisit.Delivery.Coordinates.Longitude,
				NodeInfo:  originalVisit.Delivery.NodeInfo,
			}
			fmt.Printf("DEBUG: Delivery para visita %d - VROOM sin coordenadas, usando originales: lat=%f, lon=%f\n", visitIndex, location.Latitude, location.Longitude)
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

	// Crear copias de las órdenes con el sequenceNumber asignado
	orders := make([]optimization.Order, len(originalVisit.Orders))
	for i, order := range originalVisit.Orders {
		orders[i] = order
		// Solo asignar sequenceNumber para pasos pickup, delivery o job (stepSequence > 0)
		if stepSequence > 0 {
			orders[i].SequenceNumber = &stepSequence
		}
	}

	return &optimization.OptimizedStep{
		Type:       stepType,
		VisitIndex: visitIndex,
		Location:   location,
		Orders:     orders,
		Arrival:    vroomStep.Arrival,
	}, nil
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
