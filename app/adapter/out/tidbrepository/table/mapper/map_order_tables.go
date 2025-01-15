package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderToTable(order domain.Order) table.Order {
	return table.Order{
		ID:                    order.ID,
		ReferenceID:           string(order.ReferenceID),
		OrganizationCountryID: 0, // Completar según la lógica de negocio
		CommerceID:            0, // Completar según la lógica de negocio
		ConsumerID:            0, // Completar según la lógica de negocio
		OrderStatusID:         0, // Completar según la lógica de negocio
		OrderTypeID:           0, // Completar según la lógica de negocio
		OrderType:             mapOrderTypeToTable(order.OrderType),
		OrderReferences:       mapReferencesToTable(order.References),
		DeliveryInstructions:  order.Destination.DeliveryInstructions,

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
		JSONItems:                         mapItemsToTable(order.Items),
		Visits:                            mapVisitsToTable(order.Visits),
		TransportRequirements:             mapTransportRequirementsToTable(order.TransportRequirements),
		Commerce:                          mapCommerceToTable(order.BusinessIdentifiers),
		Consumer:                          mapConsumerToTable(order.BusinessIdentifiers),
	}
}

func mapReferencesToTable(references []domain.Reference) []table.OrderReferences {
	mapped := make([]table.OrderReferences, len(references))
	for i, ref := range references {
		mapped[i] = table.OrderReferences{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return mapped
}

func mapItemsToTable(items []domain.Item) table.JSONItems {
	mapped := make(table.JSONItems, len(items))
	for i, item := range items {
		mapped[i] = table.Items{
			ReferenceID:       string(item.ReferenceID),
			LogisticCondition: item.LogisticCondition,
			QuantityNumber:    item.Quantity.QuantityNumber,
			QuantityUnit:      item.Quantity.QuantityUnit,
			JSONInsurance: table.JSONInsurance{
				UnitValue: item.Insurance.UnitValue,
				Currency:  item.Insurance.Currency,
			},
			Description: item.Description,
			JSONDimensions: table.JSONDimensions{
				Height: item.Dimensions.Height,
				Width:  item.Dimensions.Width,
				Depth:  item.Dimensions.Depth,
				Unit:   item.Dimensions.Unit,
			},
			JSONWeight: table.JSONWeight{
				WeightValue: item.Weight.Value,
				WeightUnit:  item.Weight.Unit,
			},
		}
	}

	return mapped
}

func MapPackagesToTable(packages []domain.Package) []table.Package {
	mapped := make([]table.Package, len(packages))
	for i, pkg := range packages {
		mapped[i] = table.Package{
			Lpn: pkg.Lpn,
			JSONDimensions: table.JSONDimensions{
				Height: pkg.Dimensions.Height,
				Width:  pkg.Dimensions.Width,
				Depth:  pkg.Dimensions.Depth,
				Unit:   pkg.Dimensions.Unit,
			},
			JSONWeight: table.JSONWeight{
				WeightValue: pkg.Weight.Value,
				WeightUnit:  pkg.Weight.Unit,
			},
			JSONInsurance: table.JSONInsurance{
				UnitValue: pkg.Insurance.UnitValue,
				Currency:  pkg.Insurance.Currency,
			},
		}
	}
	return mapped
}
func MapPackageToTable(pkg domain.Package) table.Package {
	return table.Package{
		ID:  pkg.ID,
		Lpn: pkg.Lpn,
		JSONDimensions: table.JSONDimensions{
			Height: pkg.Dimensions.Height,
			Width:  pkg.Dimensions.Width,
			Depth:  pkg.Dimensions.Depth,
			Unit:   pkg.Dimensions.Unit,
		},
		JSONWeight: table.JSONWeight{
			WeightValue: pkg.Weight.Value,
			WeightUnit:  pkg.Weight.Unit,
		},
		JSONInsurance: table.JSONInsurance{
			UnitValue: pkg.Insurance.UnitValue,
			Currency:  pkg.Insurance.Currency,
		},
	}
}

func mapVisitsToTable(visits []domain.Visit) []table.Visit {
	mappedVisits := make([]table.Visit, len(visits))

	for i, visit := range visits {
		mappedVisits[i] = table.Visit{
			Date:           visit.Date,
			TimeRangeStart: visit.TimeRange.StartTime,
			TimeRangeEnd:   visit.TimeRange.EndTime,
		}
	}

	return mappedVisits
}

func mapTransportRequirementsToTable(requirements []domain.Reference) table.JSONReference {
	// Convertir los requisitos a JSONReference
	var jsonReference table.JSONReference
	for _, req := range requirements {
		jsonReference = append(jsonReference, table.Reference{
			Type:  req.Type,
			Value: req.Value,
		})
	}
	return jsonReference
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
			ID:         0,
			FullName:   operator.Contact.FullName,
			Email:      operator.Contact.Email,
			Phone:      operator.Contact.Phone,
			NationalID: operator.Contact.NationalID,
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
		Province:     address.Province,
		TimeZone:     address.TimeZone,
	}
}
