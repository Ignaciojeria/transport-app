package mapper

import (
	"context"
	"fmt"
	"sort"
	"transport-app/app/adapter/out/fuegoapiclient/model"
	"transport-app/app/domain"
)

func MapUpsertRouteRequest(r domain.Route) model.UpsertRouteRequest {
	// Agrupar órdenes por secuencia, dirección y contacto
	visitGroups := groupOrdersByVisit(r.Orders)

	// Convertir grupos a visitas
	visits := make([]model.Visit, 0, len(visitGroups))
	for _, group := range visitGroups {
		visit := mapOrderGroupToVisit(group)
		visits = append(visits, visit)
	}

	// Ordenar visitas por número de secuencia
	sort.Slice(visits, func(i, j int) bool {
		return visits[i].SequenceNumber < visits[j].SequenceNumber
	})

	return model.UpsertRouteRequest{
		ReferenceID: r.ReferenceID,
		Plan: model.Plan{
			ReferenceID: r.ReferenceID,
		},
		Vehicle: model.Vehicle{
			Plate: r.Vehicle.Plate,
		},
		Geometry: model.Geometry{
			Encoding: "polyline",
			Type:     "linestring",
			Value:    "", // TODO: Implementar cálculo de geometría
		},
		Visits:    visits,
		CreatedAt: "", // TODO: Implementar timestamp de creación
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
func mapOrderGroupToVisit(group OrderGroup) model.Visit {
	// Mapear órdenes del grupo
	orders := make([]model.Order, 0, len(group.Orders))
	for _, order := range group.Orders {
		modelOrder := mapOrderToModel(order)
		orders = append(orders, modelOrder)
	}

	return model.Visit{
		Type:        "delivery",
		AddressInfo: mapAddressInfoToModel(group.AddressInfo),
		NodeInfo: model.NodeInfo{
			ReferenceID: "", // TODO: Implementar si es necesario
		},
		DeliveryInstructions: group.DeliveryInstructions,
		SequenceNumber:       group.SequenceNumber,
		Orders:               orders,
	}
}

// mapOrderToModel convierte una orden del dominio al modelo
func mapOrderToModel(order domain.Order) model.Order {
	deliveryUnits := make([]model.DeliveryUnit, 0, len(order.DeliveryUnits))
	for _, du := range order.DeliveryUnits {
		modelDU := mapDeliveryUnitToModel(du)
		deliveryUnits = append(deliveryUnits, modelDU)
	}

	return model.Order{
		ReferenceID:   string(order.ReferenceID),
		DeliveryUnits: deliveryUnits,
	}
}

// mapDeliveryUnitToModel convierte una unidad de entrega del dominio al modelo
func mapDeliveryUnitToModel(du domain.DeliveryUnit) model.DeliveryUnit {
	items := make([]model.Item, 0, len(du.Items))
	for _, item := range du.Items {
		modelItem := model.Item{
			Sku: item.Sku,
		}
		items = append(items, modelItem)
	}

	return model.DeliveryUnit{
		Lpn:   du.Lpn,
		Items: items,
	}
}

// mapAddressInfoToModel convierte AddressInfo del dominio al modelo
func mapAddressInfoToModel(addr domain.AddressInfo) model.AddressInfo {
	return model.AddressInfo{
		AddressLine1:  addr.AddressLine1,
		AddressLine2:  addr.AddressLine2,
		Contact:       mapContactToModel(addr.Contact),
		Coordinates:   mapCoordinatesToModel(addr.Coordinates),
		PoliticalArea: mapPoliticalAreaToModel(addr.PoliticalArea),
		ZipCode:       addr.ZipCode,
	}
}

// mapContactToModel convierte Contact del dominio al modelo
func mapContactToModel(contact domain.Contact) model.Contact {
	return model.Contact{
		Email:      contact.PrimaryEmail,
		FullName:   contact.FullName,
		NationalID: contact.NationalID,
		Phone:      contact.PrimaryPhone,
	}
}

// mapCoordinatesToModel convierte Coordinates del dominio al modelo
func mapCoordinatesToModel(coords domain.Coordinates) model.Coordinates {
	return model.Coordinates{
		Latitude:  coords.Point[1], // orb.Point es [longitude, latitude]
		Longitude: coords.Point[0],
	}
}

// mapPoliticalAreaToModel convierte PoliticalArea del dominio al modelo
func mapPoliticalAreaToModel(pa domain.PoliticalArea) model.PoliticalArea {
	return model.PoliticalArea{
		Code:     pa.Code,
		District: pa.District,
		Province: pa.Province,
		State:    pa.State,
	}
}
