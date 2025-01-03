package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderToTable(order domain.Order) table.Order {
	return table.Order{
		ID:                                order.ID,
		ReferenceID:                       string(order.ReferenceID),
		OrganizationCountryID:             0,
		CommerceID:                        0,
		ConsumerID:                        0,
		OrderStatusID:                     0,
		OrderTypeID:                       0,
		OrderType:                         mapOrderTypeToTable(order.OrderType),
		TransportOrderReferences:          mapReferencesToTable(order.References),
		DeliveryInstructions:              order.Destination.DeliveryInstructions,
		OriginID:                          0,
		Origin:                            mapOriginToTable(order.Origin),
		DestinationID:                     0,
		Destination:                       mapDestinationToTable(order.Destination),
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
		TransportRequirementsReferences:   mapTransportRequirementsToTable(order.TransportRequirements),
		//	OrganizationCountry:               MapOrganizationToTable(order.Organization),
		Commerce: mapCommerceToTable(order.BusinessIdentifiers),
		Consumer: mapConsumerToTable(order.BusinessIdentifiers),
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

func mapTransportRequirementsToTable(requirements []domain.References) []table.TransportRequirementsReferences {
	mapped := make([]table.TransportRequirementsReferences, len(requirements))
	for i, req := range requirements {
		mapped[i] = table.TransportRequirementsReferences{
			Type:  req.Type,
			Value: req.Value,
		}
	}
	return mapped
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
	if t.Type == "" {
		t.Type = "UNSPECIFIED"
		t.Description = "Order type not specified"
	}
	return table.OrderType{
		Type:        t.Type,
		Description: t.Description,
	}
}

func mapOriginToTable(origin domain.Origin) table.Origin {
	return table.Origin{
		NodeInfoID:    0, // This would depend on the domain logic for NodeInfo mapping
		AddressInfoID: 0, // This can be replaced with actual logic to map AddressInfo
		AddressInfo:   mapAddressInfoToTable(origin.AddressInfo),
	}
}

func mapDestinationToTable(destination domain.Destination) table.Destination {
	return table.Destination{
		NodeInfoID:    0, // This would depend on the domain logic for NodeInfo mapping
		AddressInfoID: 0, // This can be replaced with actual logic to map AddressInfo
		AddressInfo:   mapAddressInfoToTable(destination.AddressInfo),
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
		Latitude:     address.Latitude,
		Longitude:    address.Longitude,
		ZipCode:      address.ZipCode,
		TimeZone:     address.TimeZone,
	}
}
