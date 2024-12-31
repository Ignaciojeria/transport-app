package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapTransportOrderToTable(order domain.TransportOrder) table.TransportOrder {
	return table.TransportOrder{
		ID:                                order.ID,
		ReferenceID:                       string(order.ReferenceID),
		OrganizationID:                    0,
		CommerceID:                        0,
		ConsumerID:                        0,
		OrderStatusID:                     0,
		OrderTypeID:                       0,
		TransportOrderReferences:          mapReferencesToTable(order.References),
		DeliveryInstructions:              order.Destination.DeliveryInstructions,
		OriginID:                          0,
		DestinationID:                     0,
		CollectAvailabilityDate:           order.CollectAvailabilityDate.Date,
		CollectAvailabilityTimeRangeStart: order.CollectAvailabilityDate.TimeRange.Start,
		CollectAvailabilityTimeRangeEnd:   order.CollectAvailabilityDate.TimeRange.End,
		PromisedDateRangeStart:            order.PromisedDate.DateRange.StartDate,
		PromisedDateRangeEnd:              order.PromisedDate.DateRange.EndDate,
		PromisedTimeRangeStart:            order.PromisedDate.TimeRange.Start,
		PromisedTimeRangeEnd:              order.PromisedDate.TimeRange.End,
		Items:                             mapItemsToTable(order.Items),
		Packages:                          mapPackagesToTable(order.Packages),
		Visit:                             mapVisitToTable(order.Visit),
		TransportRequirementsReferences:   mapTransportRequirementsToTable(order.TransportRequirements),
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
		TimeRangeStart: visit.TimeRange.Start,
		TimeRangeEnd:   visit.TimeRange.End,
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
