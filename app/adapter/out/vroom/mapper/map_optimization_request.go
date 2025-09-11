package mapper

import (
	"context"
	"fmt"
	"transport-app/app/adapter/out/vroom/model"
	"transport-app/app/domain"
	"transport-app/app/domain/optimization"
)

type skillRegistry struct {
	counter int64
	mapping map[string]int64
}

type locationContactRegistry struct {
	counter int64
	mapping map[string]int64
}

func newSkillRegistry() *skillRegistry {
	return &skillRegistry{
		counter: 1,
		mapping: make(map[string]int64),
	}
}

func newLocationContactRegistry() *locationContactRegistry {
	return &locationContactRegistry{
		counter: 1,
		mapping: make(map[string]int64),
	}
}

func (r *skillRegistry) getSkillID(skill string) int64 {
	if id, ok := r.mapping[skill]; ok {
		return id
	}
	r.mapping[skill] = r.counter
	r.counter++
	return r.mapping[skill]
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

func generateLocationKey(lat, lon float64) string {
	return fmt.Sprintf("%.6f_%.6f", lat, lon)
}

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

// calculateVisitCapacity calcula la capacidad total de una visita basada en sus orders
func calculateVisitCapacity(visit optimization.Visit) (totalWeight, totalVolume, totalInsurance int64) {
	for _, order := range visit.Orders {
		for _, deliveryUnit := range order.DeliveryUnits {
			totalWeight += deliveryUnit.Weight
			totalVolume += deliveryUnit.Volume
			totalInsurance += deliveryUnit.Price
		}
	}
	return
}

func MapOptimizationRequest(ctx context.Context, req optimization.FleetOptimization) (model.VroomOptimizationRequest, error) {
	registry := newSkillRegistry()
	locationRegistry := newLocationContactRegistry()

	var vehicles []model.VroomVehicle
	for i, v := range req.Vehicles {
		vehicle := model.VroomVehicle{
			ID: i + 1,
		}

		// Solo incluir Start si las coordenadas no son cero
		if v.StartLocation.AddressInfo.Coordinates.Longitude != 0 || v.StartLocation.AddressInfo.Coordinates.Latitude != 0 {
			vehicle.Start = &[2]float64{
				v.StartLocation.AddressInfo.Coordinates.Longitude,
				v.StartLocation.AddressInfo.Coordinates.Latitude,
			}
		}

		// Solo incluir End si las coordenadas no son cero
		if v.EndLocation.AddressInfo.Coordinates.Longitude != 0 || v.EndLocation.AddressInfo.Coordinates.Latitude != 0 {
			vehicle.End = &[2]float64{
				v.EndLocation.AddressInfo.Coordinates.Longitude,
				v.EndLocation.AddressInfo.Coordinates.Latitude,
			}
		}

		// Incluir Capacity siempre con los 3 valores en orden: [peso, volumen, seguro]
		vehicle.Capacity = []int64{
			v.Capacity.Weight,
			v.Capacity.Volume,
			v.Capacity.Insurance,
		}

		// Solo incluir Skills si no está vacío
		if len(v.Skills) > 0 {
			vehicle.Skills = mapSkills(v.Skills, registry)
		}

		// Solo incluir TimeWindow si los valores son válidos
		if v.TimeWindow.Start != "" && v.TimeWindow.End != "" {
			vehicle.TimeWindow = parseTimeRange(v.TimeWindow.Start, v.TimeWindow.End)
		}

		vehicles = append(vehicles, vehicle)
	}

	var jobs []model.VroomJob
	var shipments []model.VroomShipment

	jobCounter := 1
	shipmentCounter := 1

	for _, visit := range req.Visits {
		// Calcular capacidad de la visita
		totalWeight, totalVolume, totalInsurance := calculateVisitCapacity(visit)

		// Verificar si hay pickup válido (coordenadas no son cero)
		hasValidPickup := visit.Pickup.AddressInfo.Coordinates.Longitude != 0 || visit.Pickup.AddressInfo.Coordinates.Latitude != 0

		// Verificar si hay delivery válido (coordenadas no son cero)
		hasValidDelivery := visit.Delivery.AddressInfo.Coordinates.Longitude != 0 || visit.Delivery.AddressInfo.Coordinates.Latitude != 0

		/* Log de depuración
		fmt.Printf("Visita: pickup=(%.6f, %.6f) delivery=(%.6f, %.6f) hasValidPickup=%v hasValidDelivery=%v orders=%d\n",
			visit.Pickup.AddressInfo.Coordinates.Longitude, visit.Pickup.AddressInfo.Coordinates.Latitude,
			visit.Delivery.AddressInfo.Coordinates.Longitude, visit.Delivery.AddressInfo.Coordinates.Latitude,
			hasValidPickup, hasValidDelivery, len(visit.Orders))
		*/
		// Si no hay delivery válido, omitir esta visita
		if !hasValidDelivery {
			//fmt.Printf("Omitiendo visita: no hay delivery válido\n")
			continue
		}

		// Si no hay pickup válido, crear un Job (entrega directa)
		if !hasValidPickup {
			//fmt.Printf("Creando Job para visita (solo delivery) con %d órdenes\n", len(visit.Orders))
			job := model.VroomJob{
				ID: jobCounter,
				Location: [2]float64{
					visit.Delivery.AddressInfo.Coordinates.Longitude,
					visit.Delivery.AddressInfo.Coordinates.Latitude,
				},
			}
			jobCounter++

			// Incluir Amount siempre con los 3 valores en orden: [peso, volumen, seguro]
			job.Amount = []int64{
				totalWeight,
				totalVolume,
				totalInsurance,
			}

			// Recopilar skills de delivery units
			var jobSkills []string
			for _, order := range visit.Orders {
				for _, deliveryUnit := range order.DeliveryUnits {
					if len(deliveryUnit.Skills) > 0 {
						jobSkills = append(jobSkills, deliveryUnit.Skills...)
					}
				}
			}
			// Solo incluir Skills si no está vacío
			if len(jobSkills) > 0 {
				job.Skills = mapSkills(jobSkills, registry)
			}

			// Solo incluir Service si no es cero
			if visit.Delivery.ServiceTime != 0 {
				job.Service = visit.Delivery.ServiceTime
			}

			// Solo incluir TimeWindows si los valores son válidos
			if visit.Delivery.TimeWindow.Start != "" && visit.Delivery.TimeWindow.End != "" {
				job.TimeWindows = [][]int{parseTimeRange(visit.Delivery.TimeWindow.Start, visit.Delivery.TimeWindow.End)}
			}

			// Solo incluir CustomUserData si hay orders o información de contacto
			customData := make(map[string]any)
			if len(visit.Orders) > 0 {
				customData["orders"] = visit.Orders
			}
			if visit.Delivery.AddressInfo.Contact.FullName != "" || visit.Delivery.AddressInfo.Contact.Email != "" || visit.Delivery.AddressInfo.Contact.Phone != "" {
				customData["delivery_contact"] = visit.Delivery.AddressInfo.Contact
			}
			if len(customData) > 0 {
				job.CustomUserData = customData
			}

			jobs = append(jobs, job)
		} else {
			// Si hay pickup válido, crear un Shipment (pickup + delivery)
			//fmt.Printf("Creando Shipment para visita (pickup + delivery) con %d órdenes\n", len(visit.Orders))
			pickupLocationKey := generateLocationKey(visit.Pickup.AddressInfo.Coordinates.Latitude, visit.Pickup.AddressInfo.Coordinates.Longitude)
			pickupContactID := getContactID(ctx, visit.Pickup.AddressInfo.Contact)
			pickupID := locationRegistry.getLocationContactID(pickupLocationKey, pickupContactID)

			deliveryLocationKey := generateLocationKey(visit.Delivery.AddressInfo.Coordinates.Latitude, visit.Delivery.AddressInfo.Coordinates.Longitude)
			deliveryContactID := getContactID(ctx, visit.Delivery.AddressInfo.Contact)
			deliveryID := locationRegistry.getLocationContactID(deliveryLocationKey, deliveryContactID)

			pickup := model.VroomStep{
				ID: int(pickupID),
				Location: &[2]float64{
					visit.Pickup.AddressInfo.Coordinates.Longitude,
					visit.Pickup.AddressInfo.Coordinates.Latitude,
				},
			}

			// Solo incluir TimeWindows si los valores son válidos
			if visit.Pickup.TimeWindow.Start != "" && visit.Pickup.TimeWindow.End != "" {
				pickup.TimeWindows = [][]int{parseTimeRange(visit.Pickup.TimeWindow.Start, visit.Pickup.TimeWindow.End)}
			}

			delivery := model.VroomStep{
				ID: int(deliveryID),
				Location: &[2]float64{
					visit.Delivery.AddressInfo.Coordinates.Longitude,
					visit.Delivery.AddressInfo.Coordinates.Latitude,
				},
			}

			// Solo incluir TimeWindows si los valores son válidos
			if visit.Delivery.TimeWindow.Start != "" && visit.Delivery.TimeWindow.End != "" {
				delivery.TimeWindows = [][]int{parseTimeRange(visit.Delivery.TimeWindow.Start, visit.Delivery.TimeWindow.End)}
			}

			shipment := model.VroomShipment{
				ID:       shipmentCounter,
				Pickup:   pickup,
				Delivery: delivery,
			}
			shipmentCounter++

			// Incluir Amount siempre con los 3 valores en orden: [peso, volumen, seguro]
			shipment.Amount = []int64{
				totalWeight,
				totalVolume,
				totalInsurance,
			}

			// Recopilar skills de delivery units
			var shipmentSkills []string
			for _, order := range visit.Orders {
				for _, deliveryUnit := range order.DeliveryUnits {
					if len(deliveryUnit.Skills) > 0 {
						shipmentSkills = append(shipmentSkills, deliveryUnit.Skills...)
					}
				}
			}
			// Solo incluir Skills si no está vacío
			if len(shipmentSkills) > 0 {
				shipment.Skills = mapSkills(shipmentSkills, registry)
			}

			// Solo incluir Service si no es cero
			if visit.Pickup.ServiceTime != 0 || visit.Delivery.ServiceTime != 0 {
				shipment.Service = visit.Pickup.ServiceTime + visit.Delivery.ServiceTime
			}

			// Solo incluir CustomUserData si hay orders o información de contacto
			customData := make(map[string]any)
			if len(visit.Orders) > 0 {
				customData["orders"] = visit.Orders
			}

			// Incluir información de contacto en CustomUserData
			if visit.Pickup.AddressInfo.Contact.FullName != "" || visit.Pickup.AddressInfo.Contact.Email != "" || visit.Pickup.AddressInfo.Contact.Phone != "" {
				customData["pickup_contact"] = visit.Pickup.AddressInfo.Contact
			}
			if visit.Delivery.AddressInfo.Contact.FullName != "" || visit.Delivery.AddressInfo.Contact.Email != "" || visit.Delivery.AddressInfo.Contact.Phone != "" {
				customData["delivery_contact"] = visit.Delivery.AddressInfo.Contact
			}

			if len(customData) > 0 {
				shipment.CustomUserData = customData
			}

			shipments = append(shipments, shipment)
		}
	}

	// Log final
	//fmt.Printf("Total jobs creados: %d, Total shipments creados: %d\n", len(jobs), len(shipments))

	return model.VroomOptimizationRequest{
		Vehicles:  vehicles,
		Jobs:      jobs,
		Shipments: shipments,
	}, nil
}

func parseTimeRange(start, end string) []int {
	// Convierte "08:00" a segundos desde medianoche
	// Retorna [inicio, fin] en segundos
	return []int{
		toSeconds(start),
		toSeconds(end),
	}
}

func toSeconds(timeStr string) int {
	// espera formato "HH:MM"
	var h, m int
	fmt.Sscanf(timeStr, "%02d:%02d", &h, &m)
	return h*3600 + m*60
}

func mapSkills(skills []string, registry *skillRegistry) []int64 {
	var mapped []int64
	for _, skill := range skills {
		mapped = append(mapped, registry.getSkillID(skill))
	}
	return mapped
}
