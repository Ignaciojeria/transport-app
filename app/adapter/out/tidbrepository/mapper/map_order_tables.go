package mapper

import (
	"encoding/json"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderToTable(order domain.Order) table.Order {
	return table.Order{
		ID:                       order.ID,
		ReferenceID:              string(order.ReferenceID),
		OrganizationCountryID:    0, // Completar según la lógica de negocio
		CommerceID:               0, // Completar según la lógica de negocio
		ConsumerID:               0, // Completar según la lógica de negocio
		OrderStatusID:            0, // Completar según la lógica de negocio
		OrderTypeID:              0, // Completar según la lógica de negocio
		OrderType:                mapOrderTypeToTable(order.OrderType),
		TransportOrderReferences: mapReferencesToTable(order.References),
		DeliveryInstructions:     order.Destination.DeliveryInstructions,

		// Origen
		OriginNodeInfoID: 0, // Completar según la lógica de negocio
		OriginNodeInfo:   mapNodeInfoToTable(order.Origin.NodeInfo),

		OriginAddressInfoID: 0, // Completar según la lógica de negocio
		OriginAddressInfo:   mapAddressInfoToTable(order.Origin.AddressInfo),

		DestinationAddressInfoID: 0, // Completar según la lógica de negocio
		DestinationAddressInfo:   mapAddressInfoToTable(order.Destination.AddressInfo),
		OriginContactID:          0,
		OriginContact:            MapContactToTable(order.Destination.AddressInfo.Contact, 0),
		DestinationContactID:     0,
		DestinationContact:       MapContactToTable(order.Destination.AddressInfo.Contact, 0),
		// Destino
		DestinationNodeInfoID: 0, // Completar según la lógica de negocio
		DestinationNodeInfo:   mapNodeInfoToTable(order.Destination.NodeInfo),

		CollectAvailabilityDate:           order.CollectAvailabilityDate.Date,
		CollectAvailabilityTimeRangeStart: order.CollectAvailabilityDate.TimeRange.StartTime,
		CollectAvailabilityTimeRangeEnd:   order.CollectAvailabilityDate.TimeRange.EndTime,
		PromisedDateRangeStart:            order.PromisedDate.DateRange.StartDate,
		PromisedDateRangeEnd:              order.PromisedDate.DateRange.EndDate,
		PromisedTimeRangeStart:            order.PromisedDate.TimeRange.StartTime,
		PromisedTimeRangeEnd:              order.PromisedDate.TimeRange.EndTime,
		Items:                             mapItemsToTable(order.Items),
		Packages:                          mapPackagesToTable(order.Packages),
		Visit:                             mapVisitToTable(order.Visit),
		TransportRequirements:             mapTransportRequirementsToTable(order.TransportRequirements),
		Commerce:                          mapCommerceToTable(order.BusinessIdentifiers),
		Consumer:                          mapConsumerToTable(order.BusinessIdentifiers),
	}
}

func mapReferencesToTable(references []domain.References) []table.TransportOrderReferences {
	mapped := make([]table.TransportOrderReferences, len(references))
	for i, ref := range references {
		mapped[i] = table.TransportOrderReferences{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return mapped
}

func mapItemsToTable(items []domain.Items) []table.Items {
	mapped := make([]table.Items, len(items))
	for i, item := range items {
		mapped[i] = table.Items{
			ReferenceID:       string(item.ReferenceID),
			LogisticCondition: item.LogisticCondition,
			Quantity: table.Quantity{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
			Insurance: table.Insurance{
				UnitValue: item.Insurance.UnitValue,
				Currency:  item.Insurance.Currency,
			},
			Description: item.Description,
			Dimensions: table.Dimensions{
				Height: item.Dimensions.Height,
				Width:  item.Dimensions.Width,
				Depth:  item.Dimensions.Depth,
				Unit:   item.Dimensions.Unit,
			},
			Weight: table.Weight{
				Value: item.Weight.Value,
				Unit:  item.Weight.Unit,
			},
		}
	}
	return mapped
}

func mapPackagesToTable(packages []domain.Packages) []table.Packages {
	mapped := make([]table.Packages, len(packages))
	for i, pkg := range packages {
		mapped[i] = table.Packages{
			Lpn:         pkg.Lpn,
			PackageType: pkg.PackageType,
			Dimensions: table.Dimensions{
				Height: pkg.Dimensions.Height,
				Width:  pkg.Dimensions.Width,
				Depth:  pkg.Dimensions.Depth,
				Unit:   pkg.Dimensions.Unit,
			},
			Weight: table.Weight{
				Value: pkg.Weight.Value,
				Unit:  pkg.Weight.Unit,
			},
			Insurance: table.Insurance{
				UnitValue: pkg.Insurance.UnitValue,
				Currency:  pkg.Insurance.Currency,
			},
			ItemReferences: mapItemReferencesToTable(pkg.ItemReferences),
		}
	}
	return mapped
}

func mapItemReferencesToTable(references []domain.ItemReferences) []table.ItemReferences {
	mapped := make([]table.ItemReferences, len(references))
	for i, ref := range references {
		mapped[i] = table.ItemReferences{
			ReferenceID: string(ref.ReferenceID),
			Quantity: table.Quantity{
				QuantityNumber: ref.Quantity.QuantityNumber,
				QuantityUnit:   ref.Quantity.QuantityUnit,
			},
		}
	}
	return mapped
}

func mapVisitToTable(visit domain.Visit) table.Visit {
	return table.Visit{
		Date:           visit.Date,
		TimeRangeStart: visit.TimeRange.StartTime,
		TimeRangeEnd:   visit.TimeRange.EndTime,
	}
}

func mapTransportRequirementsToTable(requirements []domain.References) []byte {
	// Serializar los requisitos en JSON
	serialized, _ := json.Marshal(requirements)
	return serialized
}

func mapCommerceToTable(bi domain.BusinessIdentifiers) table.Commerce {
	return table.Commerce{
		Name: bi.Commerce,
	}
}

func mapConsumerToTable(bi domain.BusinessIdentifiers) table.Consumer {
	return table.Consumer{
		Name: bi.Consumer,
	}
}

func mapOrderTypeToTable(t domain.OrderType) table.OrderType {
	return table.OrderType{
		Type:        t.Type,
		Description: t.Description,
	}
}

func mapNodeInfoToTable(node domain.NodeInfo) table.NodeInfo {
	return table.NodeInfo{
		ReferenceID: string(node.ReferenceID),
		Name:        node.Name,
		Type:        node.Type,
		Operator:    mapOperatorToTable(node.Operator),
	}
}

func mapOperatorToTable(operator domain.Operator) table.Operator {
	return table.Operator{
		ID:   0,
		Type: operator.Type,
		Contact: table.Contact{
			ID:       0,
			FullName: operator.Contact.FullName,
			Email:    operator.Contact.Email,
			Phone:    operator.Contact.Phone,
		},
	}
}

func mapAddressInfoToTable(address domain.AddressInfo) table.AddressInfo {
	return table.AddressInfo{
		State:        address.State,
		County:       address.County,
		District:     address.District,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		AddressLine3: address.AddressLine3,
		RawAddress:   address.RawAddress(),
		Latitude:     address.Latitude,
		Longitude:    address.Longitude,
		ZipCode:      address.ZipCode,
		TimeZone:     address.TimeZone,
	}
}
