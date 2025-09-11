package mapper

import (
	"context"
	"fmt"
	"sort"
	"time"
	"transport-app/app/adapter/in/fuegoapi/request"
	"transport-app/app/domain"
)

func MapUpsertRouteRequest(r domain.Route) request.UpsertRouteRequest {
	// Agrupar órdenes por secuencia, dirección y contacto
	visitGroups := groupOrdersByVisit(r.Orders)

	// Convertir grupos a visitas
	visits := make([]request.UpsertRouteVisit, 0, len(visitGroups))
	for _, group := range visitGroups {
		visit := mapOrderGroupToVisit(group)
		visits = append(visits, visit)
	}

	// Ordenar visitas por número de secuencia
	sort.Slice(visits, func(i, j int) bool {
		return visits[i].SequenceNumber < visits[j].SequenceNumber
	})

	return request.UpsertRouteRequest{
		ReferenceID:     r.ReferenceID,
		PlanReferenceID: r.ReferenceID, // Usar el mismo ID como referencia del plan
		Vehicle:         mapVehicleToRequest(r.Vehicle),
		Geometry: request.UpsertRouteGeometry{
			Encoding: r.Geometry.Encoding,
			Type:     r.Geometry.Type,
			Value:    r.Geometry.Value,
		},
		Visits:    visits,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// OrderGroup representa un grupo de órdenes que se pueden agrupar en una visita
type OrderGroup struct {
	SequenceNumber       int
	AddressInfo          domain.AddressInfo
	Contact              domain.Contact
	DeliveryInstructions string
	Orders               []domain.Order
}

// groupOrdersByVisit agrupa las órdenes por secuencia, dirección y contacto
func groupOrdersByVisit(orders []domain.Order) []OrderGroup {
	if len(orders) == 0 {
		return []OrderGroup{}
	}

	// Crear un mapa para agrupar órdenes
	groups := make(map[string]OrderGroup)
	ctx := context.Background()

	for _, order := range orders {
		// Obtener secuencia (si no tiene, usar 0)
		sequence := 0
		if order.SequenceNumber != nil {
			sequence = *order.SequenceNumber
		}

		// Crear clave única basada en secuencia, dirección y contacto
		addressKey := order.Destination.AddressInfo.DocID(ctx).String()
		contactKey := order.Destination.AddressInfo.Contact.DocID(ctx).String()
		groupKey := fmt.Sprintf("%d_%s_%s", sequence, addressKey, contactKey)

		if group, exists := groups[groupKey]; exists {
			// Agregar orden al grupo existente
			group.Orders = append(group.Orders, order)
			groups[groupKey] = group
		} else {
			// Crear nuevo grupo
			groups[groupKey] = OrderGroup{
				SequenceNumber:       sequence,
				AddressInfo:          order.Destination.AddressInfo,
				Contact:              order.Destination.AddressInfo.Contact,
				DeliveryInstructions: order.DeliveryInstructions,
				Orders:               []domain.Order{order},
			}
		}
	}

	// Convertir mapa a slice
	result := make([]OrderGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}

// mapOrderGroupToVisit convierte un grupo de órdenes a una visita
func mapOrderGroupToVisit(group OrderGroup) request.UpsertRouteVisit {
	// Mapear órdenes del grupo
	orders := make([]request.UpsertRouteOrder, 0, len(group.Orders))
	for _, order := range group.Orders {
		modelOrder := mapOrderToRequest(order)
		orders = append(orders, modelOrder)
	}

	return request.UpsertRouteVisit{
		Type:           "delivery",
		AddressInfo:    mapAddressInfoToRequest(group.AddressInfo),
		NodeInfo:       mapNodeInfoToRequest(group.AddressInfo),
		SequenceNumber: group.SequenceNumber,
		ServiceTime:    0, // TODO: Implementar si es necesario
		TimeWindow: request.UpsertRouteTimeWindow{
			Start: "", // TODO: Implementar si es necesario
			End:   "", // TODO: Implementar si es necesario
		},
		Orders: orders,
	}
}

// mapOrderToRequest convierte una orden del dominio al request
func mapOrderToRequest(order domain.Order) request.UpsertRouteOrder {
	deliveryUnits := make([]request.UpsertRouteDeliveryUnit, 0, len(order.DeliveryUnits))
	for _, du := range order.DeliveryUnits {
		modelDU := mapDeliveryUnitToRequest(du)
		deliveryUnits = append(deliveryUnits, modelDU)
	}

	return request.UpsertRouteOrder{
		ReferenceID:          order.ReferenceID.String(),
		Contact:              mapContactToRequest(order.Destination.AddressInfo.Contact),
		DeliveryInstructions: order.DeliveryInstructions,
		DeliveryUnits:        deliveryUnits,
	}
}

// mapDeliveryUnitToRequest convierte una unidad de entrega del dominio al request
func mapDeliveryUnitToRequest(du domain.DeliveryUnit) request.UpsertRouteDeliveryUnit {
	items := make([]request.UpsertRouteItem, 0, len(du.Items))
	for _, item := range du.Items {
		modelItem := request.UpsertRouteItem{
			Sku:         item.Sku,
			Description: item.Description,
			Quantity:    item.Quantity,
		}
		items = append(items, modelItem)
	}

	// Manejar punteros para Volume, Weight, Insurance
	var volume, weight, insurance int64
	if du.Volume != nil {
		volume = *du.Volume
	}
	if du.Weight != nil {
		weight = *du.Weight
	}
	if du.Price != nil {
		insurance = *du.Price
	}

	// Convertir Skills de domain.Skill a string
	skills := make([]string, 0, len(du.Skills))
	for _, skill := range du.Skills {
		skills = append(skills, string(skill)) // Skill es un tipo string
	}

	return request.UpsertRouteDeliveryUnit{
		Items:  items,
		Volume: volume,
		Weight: weight,
		Price:  insurance,
		Lpn:    du.Lpn,
		Skills: skills,
	}
}

// mapVehicleToRequest convierte el vehículo del dominio al request
func mapVehicleToRequest(vehicle domain.Vehicle) request.UpsertRouteVehicle {
	// Como el dominio Vehicle no tiene StartLocation, EndLocation, Skills, etc.
	// Usamos valores por defecto o vacíos
	return request.UpsertRouteVehicle{
		Plate: vehicle.Plate,
		StartLocation: request.UpsertRouteVehicleLocation{
			AddressInfo: request.UpsertRouteAddressInfo{}, // Vacío por defecto
			NodeInfo:    request.UpsertRouteNodeInfo{},    // Vacío por defecto
		},
		EndLocation: request.UpsertRouteVehicleLocation{
			AddressInfo: request.UpsertRouteAddressInfo{}, // Vacío por defecto
			NodeInfo:    request.UpsertRouteNodeInfo{},    // Vacío por defecto
		},
		Skills: []string{}, // Vacío por defecto
		TimeWindow: request.UpsertRouteTimeWindow{
			Start: "", // TODO: Implementar si es necesario
			End:   "", // TODO: Implementar si es necesario
		},
		Capacity: request.UpsertRouteVehicleCapacity{
			Volume:                int64(vehicle.Weight.Value), // Usar Weight.Value como volumen
			Weight:                int64(vehicle.Weight.Value),
			Insurance:             int64(vehicle.Insurance.MaxInsuranceCoverage.Amount),
			DeliveryUnitsQuantity: 0, // TODO: Implementar si es necesario
		},
	}
}

// mapAddressInfoToRequest convierte AddressInfo del dominio al request
func mapAddressInfoToRequest(addr domain.AddressInfo) request.UpsertRouteAddressInfo {
	return request.UpsertRouteAddressInfo{
		AddressLine1:  addr.AddressLine1,
		AddressLine2:  addr.AddressLine2,
		Coordinates:   mapCoordinatesToRequest(addr.Coordinates),
		PoliticalArea: mapPoliticalAreaToRequest(addr.PoliticalArea),
		ZipCode:       addr.ZipCode,
	}
}

// mapContactToRequest convierte Contact del dominio al request
func mapContactToRequest(contact domain.Contact) request.UpsertRouteContact {
	return request.UpsertRouteContact{
		Email:      contact.PrimaryEmail,
		FullName:   contact.FullName,
		NationalID: contact.NationalID,
		Phone:      contact.PrimaryPhone,
	}
}

// mapCoordinatesToRequest convierte Coordinates del dominio al request
func mapCoordinatesToRequest(coords domain.Coordinates) request.UpsertRouteCoordinates {
	return request.UpsertRouteCoordinates{
		Latitude:  coords.Point[1], // orb.Point es [longitude, latitude]
		Longitude: coords.Point[0],
	}
}

// mapPoliticalAreaToRequest convierte PoliticalArea del dominio al request
func mapPoliticalAreaToRequest(pa domain.PoliticalArea) request.UpsertRoutePoliticalArea {
	return request.UpsertRoutePoliticalArea{
		Code:            pa.Code,
		AdminAreaLevel1: pa.AdminAreaLevel1,
		AdminAreaLevel2: pa.AdminAreaLevel2,
		AdminAreaLevel3: pa.AdminAreaLevel3,
		AdminAreaLevel4: pa.AdminAreaLevel4,
	}
}

// mapNodeInfoToRequest convierte la información del nodo del dominio al request
func mapNodeInfoToRequest(addr domain.AddressInfo) request.UpsertRouteNodeInfo {
	// Por ahora usar un ID vacío, se puede implementar lógica específica si es necesario
	return request.UpsertRouteNodeInfo{
		ReferenceID: "", // TODO: Implementar si es necesario
	}
}
