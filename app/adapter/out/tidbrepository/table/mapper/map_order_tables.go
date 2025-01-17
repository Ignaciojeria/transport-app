package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderToTable(order domain.Order) table.Order {
	orgCountryID := order.Organization.OrganizationCountryID
	return table.Order{
		ID:                    order.ID,
		ReferenceID:           string(order.ReferenceID),
		OrganizationCountryID: order.Organization.OrganizationCountryID, // Completar según la lógica de negocio
		CommerceID:            order.BusinessIdentifiers.CommerceID,     // Completar según la lógica de negocio
		ConsumerID:            order.BusinessIdentifiers.ConsumerID,     // Completar según la lógica de negocio
		OrderStatusID:         order.OrderStatus.ID,                     // Completar según la lógica de negocio
		OrderTypeID:           order.OrderType.ID,                       // Completar según la lógica de negocio
		OrderType:             mapOrderTypeToTable(order.OrderType, orgCountryID),
		OrderReferences:       mapReferencesToTable(order.References),
		DeliveryInstructions:  order.Destination.DeliveryInstructions,

		// Origen
		OriginNodeInfoID: order.Origin.NodeInfo.ID, // Completar según la lógica de negocio
		OriginNodeInfo:   mapNodeInfoToTable(order.Origin.NodeInfo, orgCountryID),

		OriginAddressInfoID: order.Origin.AddressInfo.ID, // Completar según la lógica de negocio
		OriginAddressInfo:   mapAddressInfoToTable(order.Origin.AddressInfo, orgCountryID),

		DestinationAddressInfoID: order.Destination.AddressInfo.ID, // Completar según la lógica de negocio
		DestinationAddressInfo:   mapAddressInfoToTable(order.Destination.AddressInfo, orgCountryID),
		OriginContactID:          order.Origin.AddressInfo.Contact.ID,
		OriginContact:            MapContactToTable(order.Destination.AddressInfo.Contact, orgCountryID),
		DestinationContactID:     order.Destination.AddressInfo.Contact.ID,
		DestinationContact:       MapContactToTable(order.Destination.AddressInfo.Contact, orgCountryID),
		// Destino
		DestinationNodeInfoID: order.Destination.NodeInfo.ID, // Completar según la lógica de negocio
		DestinationNodeInfo:   mapNodeInfoToTable(order.Destination.NodeInfo, orgCountryID),

		CollectAvailabilityDate:           order.CollectAvailabilityDate.Date,
		CollectAvailabilityTimeRangeStart: order.CollectAvailabilityDate.TimeRange.StartTime,
		CollectAvailabilityTimeRangeEnd:   order.CollectAvailabilityDate.TimeRange.EndTime,
		PromisedDateRangeStart:            order.PromisedDate.DateRange.StartDate,
		PromisedDateRangeEnd:              order.PromisedDate.DateRange.EndDate,
		PromisedTimeRangeStart:            order.PromisedDate.TimeRange.StartTime,
		PromisedTimeRangeEnd:              order.PromisedDate.TimeRange.EndTime,
		JSONItems:                         mapItemsToTable(order.Items),
		//Visits:                            mapVisitsToTable(order.Visits),
		TransportRequirements: mapTransportRequirementsToTable(order.TransportRequirements),
		Commerce:              mapCommerceToTable(order.BusinessIdentifiers, orgCountryID),
		Consumer:              mapConsumerToTable(order.BusinessIdentifiers, orgCountryID),
		Packages:              MapPackagesToTable(order.Packages, orgCountryID),
	}
}

func mapReferencesToTable(references []domain.Reference) []table.OrderReferences {
	mapped := make([]table.OrderReferences, len(references))
	for i, ref := range references {
		mapped[i] = table.OrderReferences{
			ID:    ref.ID,
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

func MapPackagesToTable(packages []domain.Package, orgCountryID int64) []table.Package {
	mapped := make([]table.Package, len(packages))
	for i, pkg := range packages {
		mapped[i] = table.Package{
			OrganizationCountryID: orgCountryID,
			ID:                    pkg.ID,
			Lpn:                   pkg.Lpn,
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
			JSONItemsReferences: mapDomainItemsToTable(pkg.ItemReferences),
		}
	}
	return mapped
}

func mapDomainItemsToTable(items []domain.ItemReference) table.JSONItemReferences {
	mapped := make(table.JSONItemReferences, len(items))
	for i, item := range items {
		mapped[i] = table.ItemReference{ // Cambiado de JSONItemReferences a JSONItemReference
			ReferenceID: string(item.ReferenceID),
			Quantity: table.Quantity{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
		}
	}
	return mapped
}

func MapPackageToTable(pkg domain.Package, orgCountryID int64) table.Package {
	return table.Package{
		OrganizationCountryID: orgCountryID,
		ID:                    pkg.ID,
		Lpn:                   pkg.Lpn,
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
		JSONItemsReferences: mapDomainItemsToTable(pkg.ItemReferences),
	}
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

func mapCommerceToTable(bi domain.BusinessIdentifiers, orgCountry int64) table.Commerce {
	return table.Commerce{
		OrganizationCountryID: orgCountry,
		ID:                    bi.CommerceID,
		Name:                  bi.Commerce,
	}
}

func mapConsumerToTable(bi domain.BusinessIdentifiers, orgCountry int64) table.Consumer {
	return table.Consumer{
		OrganizationCountryID: orgCountry,
		ID:                    bi.ConsumerID,
		Name:                  bi.Consumer,
	}
}

func mapOrderTypeToTable(t domain.OrderType, orgCountry int64) table.OrderType {
	return table.OrderType{
		OrganizationCountryID: orgCountry,
		ID:                    t.ID,
		Type:                  t.Type,
		Description:           t.Description,
	}
}

func mapNodeInfoToTable(node domain.NodeInfo, orgCountry int64) table.NodeInfo {
	return table.NodeInfo{
		ID:                    node.ID,
		OrganizationCountryID: orgCountry,
		ReferenceID:           string(node.ReferenceID),
		Name:                  node.Name,
		Type:                  node.Type,
		//Operator:              mapOperatorToTable(node.Operator, orgCountry),
	}
}

/*
func mapOperatorToTable(operator domain.Operator, orgCountry int64) table.Operator {
	return table.Operator{
		OrganizationCountryID: orgCountry,
		ID:                    operator.ID,
		Type:                  operator.Type,
		Contact: table.Contact{
			ID:         0,
			FullName:   operator.Contact.FullName,
			Email:      operator.Contact.Email,
			Phone:      operator.Contact.Phone,
			NationalID: operator.Contact.NationalID,
		},
	}
}*/

func mapAddressInfoToTable(address domain.AddressInfo, orgCountry int64) table.AddressInfo {
	return table.AddressInfo{
		OrganizationCountryID: orgCountry,
		ID:                    address.ID,
		State:                 address.State,
		County:                address.County,
		District:              address.District,
		AddressLine1:          address.AddressLine1,
		AddressLine2:          address.AddressLine2,
		AddressLine3:          address.AddressLine3,
		RawAddress:            address.RawAddress(),
		Latitude:              address.Latitude,
		Longitude:             address.Longitude,
		ZipCode:               address.ZipCode,
		Province:              address.Province,
		TimeZone:              address.TimeZone,
	}
}
