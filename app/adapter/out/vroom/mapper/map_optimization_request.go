package mapper

import (
	"context"
	"fmt"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/adapter/out/vroom/model"
	"transport-app/app/domain"
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

func getContactID(ctx context.Context, contact struct {
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	NationalID string `json:"nationalID"`
	FullName   string `json:"fullName"`
}) string {
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
func calculateVisitCapacity(visit struct {
	Pickup struct {
		Coordinates struct {
			Latitude  float64 `json:"latitude" example:"-33.45" description:"Pickup point latitude"`
			Longitude float64 `json:"longitude" example:"-70.66" description:"Pickup point longitude"`
		} `json:"coordinates"`
		ServiceTime int64 `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
		Contact     struct {
			Email      string `json:"email"`
			Phone      string `json:"phone"`
			NationalID string `json:"nationalID"`
			FullName   string `json:"fullName"`
		} `json:"contact"`
		Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
		TimeWindow struct {
			Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
			End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
		} `json:"timeWindow"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"pickup"`
	Delivery struct {
		Coordinates struct {
			Latitude  float64 `json:"latitude" example:"-33.45" description:"Pickup point latitude"`
			Longitude float64 `json:"longitude" example:"-70.66" description:"Pickup point longitude"`
		} `json:"coordinates"`
		ServiceTime int64 `json:"serviceTime" example:"30" description:"Time in seconds required to complete the service at this location"`
		Contact     struct {
			Email      string `json:"email"`
			Phone      string `json:"phone"`
			NationalID string `json:"nationalID"`
			FullName   string `json:"fullName"`
		} `json:"contact"`
		Skills     []string `json:"skills" description:"Required vehicle capabilities for this visit"`
		TimeWindow struct {
			Start string `json:"start" example:"09:00" description:"Visit time window start (24h format)"`
			End   string `json:"end" example:"17:00" description:"Visit time window end (24h format)"`
		} `json:"timeWindow"`
		NodeInfo struct {
			ReferenceID string `json:"referenceID"`
		} `json:"nodeInfo"`
	} `json:"delivery"`
	Orders []struct {
		DeliveryUnits []struct {
			Items []struct {
				Sku string `json:"sku" example:"SKU123" description:"Stock keeping unit identifier"`
			} `json:"items"`
			Insurance int64  `json:"insurance" example:"10000" description:"Insurance value of the delivery unit"`
			Volume    int64  `json:"volume" example:"1000" description:"Volume of the delivery unit in cubic meters"`
			Weight    int64  `json:"weight" example:"1000" description:"Weight of the delivery unit in grams"`
			Lpn       string `json:"lpn" example:"LPN456" description:"License plate number of the delivery unit"`
		} `json:"deliveryUnits"`
		ReferenceID string `json:"referenceID" example:"ORD789" description:"Unique identifier for the order"`
	} `json:"orders"`
}) (totalWeight, totalDeliveryUnits, totalInsurance int64) {
	for _, order := range visit.Orders {
		for _, deliveryUnit := range order.DeliveryUnits {
			totalWeight += deliveryUnit.Weight
			totalInsurance += deliveryUnit.Insurance
			totalDeliveryUnits++
		}
	}
	return
}

func MapOptimizationRequest(ctx context.Context, req request.OptimizationRequest) (model.VroomOptimizationRequest, error) {
	registry := newSkillRegistry()
	locationRegistry := newLocationContactRegistry()

	var vehicles []model.VroomVehicle
	for i, v := range req.Vehicles {
		vehicle := model.VroomVehicle{
			ID: i + 1,
		}

		// Solo incluir Start si las coordenadas no son cero
		if v.StartLocation.Longitude != 0 || v.StartLocation.Latitude != 0 {
			vehicle.Start = &[2]float64{
				v.StartLocation.Longitude,
				v.StartLocation.Latitude,
			}
		}

		// Solo incluir End si las coordenadas no son cero
		if v.EndLocation.Longitude != 0 || v.EndLocation.Latitude != 0 {
			vehicle.End = &[2]float64{
				v.EndLocation.Longitude,
				v.EndLocation.Latitude,
			}
		}

		// Solo incluir Capacity si al menos un valor no es cero
		if v.Capacity.Weight != 0 || v.Capacity.DeliveryUnitsQuantity != 0 || v.Capacity.Insurance != 0 {
			vehicle.Capacity = []int64{
				v.Capacity.Weight,
				v.Capacity.DeliveryUnitsQuantity,
				v.Capacity.Insurance,
			}
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

	for i, visit := range req.Visits {
		// Verificar si tenemos pickup válido
		hasValidPickup := visit.Pickup.Coordinates.Longitude != 0 || visit.Pickup.Coordinates.Latitude != 0
		// Verificar si tenemos delivery válido
		hasValidDelivery := visit.Delivery.Coordinates.Longitude != 0 || visit.Delivery.Coordinates.Latitude != 0

		// Calcular capacidad de la visita
		totalWeight, totalDeliveryUnits, totalInsurance := calculateVisitCapacity(visit)

		if hasValidPickup && hasValidDelivery {
			// Ambos son válidos -> crear Shipment

			// Generar identificadores únicos para pickup y delivery basados en coordenadas y contacto
			pickupLocationKey := generateLocationKey(visit.Pickup.Coordinates.Latitude, visit.Pickup.Coordinates.Longitude)
			pickupContactID := getContactID(ctx, visit.Pickup.Contact)
			pickupID := locationRegistry.getLocationContactID(pickupLocationKey, pickupContactID)

			deliveryLocationKey := generateLocationKey(visit.Delivery.Coordinates.Latitude, visit.Delivery.Coordinates.Longitude)
			deliveryContactID := getContactID(ctx, visit.Delivery.Contact)
			deliveryID := locationRegistry.getLocationContactID(deliveryLocationKey, deliveryContactID)

			pickup := model.VroomStep{
				ID: int(pickupID),
			}

			// Solo incluir Location en pickup si las coordenadas no son cero
			if visit.Pickup.Coordinates.Longitude != 0 || visit.Pickup.Coordinates.Latitude != 0 {
				pickup.Location = &[2]float64{
					visit.Pickup.Coordinates.Longitude,
					visit.Pickup.Coordinates.Latitude,
				}
			}

			// Solo incluir TimeWindows si los valores son válidos
			if visit.Pickup.TimeWindow.Start != "" && visit.Pickup.TimeWindow.End != "" {
				pickup.TimeWindows = [][]int{parseTimeRange(visit.Pickup.TimeWindow.Start, visit.Pickup.TimeWindow.End)}
			}

			delivery := model.VroomStep{
				ID: int(deliveryID),
			}

			// Solo incluir Location en delivery si las coordenadas no son cero
			if visit.Delivery.Coordinates.Longitude != 0 || visit.Delivery.Coordinates.Latitude != 0 {
				delivery.Location = &[2]float64{
					visit.Delivery.Coordinates.Longitude,
					visit.Delivery.Coordinates.Latitude,
				}
			}

			// Solo incluir TimeWindows si los valores son válidos
			if visit.Delivery.TimeWindow.Start != "" && visit.Delivery.TimeWindow.End != "" {
				delivery.TimeWindows = [][]int{parseTimeRange(visit.Delivery.TimeWindow.Start, visit.Delivery.TimeWindow.End)}
			}

			shipment := model.VroomShipment{
				ID:       i + 1,
				Pickup:   pickup,
				Delivery: delivery,
			}

			// Solo incluir Amount si al menos un valor no es cero
			if totalWeight != 0 || totalInsurance != 0 {
				shipment.Amount = []int64{
					totalWeight,
					totalDeliveryUnits,
					totalInsurance,
				}
			}

			// Solo incluir Skills si no está vacío (usar skills de pickup o delivery)
			var skills []string
			if len(visit.Pickup.Skills) > 0 {
				skills = append(skills, visit.Pickup.Skills...)
			}
			if len(visit.Delivery.Skills) > 0 {
				skills = append(skills, visit.Delivery.Skills...)
			}
			if len(skills) > 0 {
				shipment.Skills = mapSkills(skills, registry)
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
			if visit.Pickup.Contact.FullName != "" || visit.Pickup.Contact.Email != "" || visit.Pickup.Contact.Phone != "" {
				customData["pickup_contact"] = visit.Pickup.Contact
			}
			if visit.Delivery.Contact.FullName != "" || visit.Delivery.Contact.Email != "" || visit.Delivery.Contact.Phone != "" {
				customData["delivery_contact"] = visit.Delivery.Contact
			}

			if len(customData) > 0 {
				shipment.CustomUserData = customData
			}

			shipments = append(shipments, shipment)

		} else if hasValidDelivery {
			// Solo delivery válido -> crear Job

			// Generar identificador único para delivery basado en coordenadas y contacto
			deliveryLocationKey := generateLocationKey(visit.Delivery.Coordinates.Latitude, visit.Delivery.Coordinates.Longitude)
			deliveryContactID := getContactID(ctx, visit.Delivery.Contact)
			deliveryID := locationRegistry.getLocationContactID(deliveryLocationKey, deliveryContactID)

			job := model.VroomJob{
				ID: int(deliveryID),
				Location: [2]float64{
					visit.Delivery.Coordinates.Longitude,
					visit.Delivery.Coordinates.Latitude,
				},
			}

			// Solo incluir TimeWindows si los valores son válidos
			if visit.Delivery.TimeWindow.Start != "" && visit.Delivery.TimeWindow.End != "" {
				job.TimeWindows = [][]int{parseTimeRange(visit.Delivery.TimeWindow.Start, visit.Delivery.TimeWindow.End)}
			}

			// Solo incluir Amount si al menos un valor no es cero
			if totalWeight != 0 || totalInsurance != 0 {
				job.Amount = []int64{
					totalWeight,
					totalDeliveryUnits,
					totalInsurance,
				}
			}

			// Solo incluir Skills si no está vacío
			if len(visit.Delivery.Skills) > 0 {
				job.Skills = mapSkills(visit.Delivery.Skills, registry)
			}

			// Solo incluir Service si no es cero
			if visit.Delivery.ServiceTime != 0 {
				job.Service = visit.Delivery.ServiceTime
			}

			// Solo incluir CustomUserData si hay orders o información de contacto
			customData := make(map[string]any)
			if len(visit.Orders) > 0 {
				customData["orders"] = visit.Orders
			}

			// Incluir información de contacto en CustomUserData
			if visit.Delivery.Contact.FullName != "" || visit.Delivery.Contact.Email != "" || visit.Delivery.Contact.Phone != "" {
				customData["delivery_contact"] = visit.Delivery.Contact
			}

			if len(customData) > 0 {
				job.CustomUserData = customData
			}

			jobs = append(jobs, job)
		}
		// Si no hay ni pickup ni delivery válidos, se omite la visita
	}

	return model.VroomOptimizationRequest{
		Vehicles:  vehicles,
		Jobs:      jobs,
		Shipments: shipments,
		Options: &model.VroomOptions{
			G:                true,
			Steps:            true,
			Overview:         true,
			MinimizeVehicles: true,
		},
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
