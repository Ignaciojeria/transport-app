package mapper

import (
	"transport-app/app/adapter/in/graphql/graph/model"
	"transport-app/app/domain"
)

// MapDeliveryUnitsFilter convierte el filtro de GraphQL al filtro del dominio
func MapDeliveryUnitsFilter(filter *model.DeliveryUnitsReportFilterInput) *domain.DeliveryUnitsFilter {
	if filter == nil {
		return nil
	}

	deliveryUnitsFilter := &domain.DeliveryUnitsFilter{}

	if filter.OnlyLatestStatus != nil {
		deliveryUnitsFilter.OnlyLatestStatus = *filter.OnlyLatestStatus
	}

	// Filtros de Order
	if filter.Order != nil {
		orderFilter := &domain.OrderFilter{}

		if len(filter.Order.ReferenceIds) > 0 {
			referenceIds := make([]string, len(filter.Order.ReferenceIds))
			for i, ref := range filter.Order.ReferenceIds {
				if ref != nil {
					referenceIds[i] = *ref
				}
			}
			orderFilter.ReferenceIds = referenceIds
		}

		if len(filter.Order.References) > 0 {
			references := make([]domain.ReferenceFilter, len(filter.Order.References))
			for i, ref := range filter.Order.References {
				if ref != nil {
					references[i] = domain.ReferenceFilter{
						Type:  ref.Type,
						Value: ref.Value,
					}
				}
			}
			orderFilter.References = references
		}

		if filter.Order.OrderType != nil {
			orderFilter.OrderType = &domain.OrderTypeFilter{
				Type:        *filter.Order.OrderType.Type,
				Description: *filter.Order.OrderType.Description,
			}
		}

		if filter.Order.GroupBy != nil {
			orderFilter.GroupBy = &domain.GroupByFilter{
				Type:  *filter.Order.GroupBy.Type,
				Value: *filter.Order.GroupBy.Value,
			}
		}

		deliveryUnitsFilter.Order = orderFilter
	}

	// Filtros de DeliveryUnit
	if filter.DeliveryUnit != nil {
		deliveryUnitFilter := &domain.DeliveryUnitFilter{}

		if len(filter.DeliveryUnit.Lpns) > 0 {
			lpns := make([]string, len(filter.DeliveryUnit.Lpns))
			for i, lpn := range filter.DeliveryUnit.Lpns {
				if lpn != nil {
					lpns[i] = *lpn
				}
			}
			deliveryUnitFilter.Lpns = lpns
		}

		if len(filter.DeliveryUnit.SizeCategories) > 0 {
			sizeCategories := make([]string, len(filter.DeliveryUnit.SizeCategories))
			for i, size := range filter.DeliveryUnit.SizeCategories {
				if size != nil {
					sizeCategories[i] = *size
				}
			}
			deliveryUnitFilter.SizeCategories = sizeCategories
		}

		if len(filter.DeliveryUnit.Labels) > 0 {
			labels := make([]domain.LabelFilter, len(filter.DeliveryUnit.Labels))
			for i, label := range filter.DeliveryUnit.Labels {
				if label != nil {
					labels[i] = domain.LabelFilter{
						Type:  label.Type,
						Value: label.Value,
					}
				}
			}
			deliveryUnitFilter.Labels = labels
		}

		deliveryUnitsFilter.DeliveryUnit = deliveryUnitFilter
	}

	// Filtros de Origin
	if filter.Origin != nil {
		originFilter := &domain.LocationFilter{}

		if len(filter.Origin.NodeReferences) > 0 {
			nodeReferences := make([]string, len(filter.Origin.NodeReferences))
			for i, ref := range filter.Origin.NodeReferences {
				if ref != nil {
					nodeReferences[i] = *ref
				}
			}
			originFilter.NodeReferences = nodeReferences
		}

		if len(filter.Origin.AddressLines) > 0 {
			addressLines := make([]string, len(filter.Origin.AddressLines))
			for i, line := range filter.Origin.AddressLines {
				if line != nil {
					addressLines[i] = *line
				}
			}
			originFilter.AddressLines = addressLines
		}

		if len(filter.Origin.AdminAreaLevel1) > 0 {
			adminAreaLevel1 := make([]string, len(filter.Origin.AdminAreaLevel1))
			for i, val := range filter.Origin.AdminAreaLevel1 {
				if val != nil {
					adminAreaLevel1[i] = *val
				}
			}
			originFilter.AdminAreaLevel1 = adminAreaLevel1
		}
		if len(filter.Origin.AdminAreaLevel2) > 0 {
			adminAreaLevel2 := make([]string, len(filter.Origin.AdminAreaLevel2))
			for i, val := range filter.Origin.AdminAreaLevel2 {
				if val != nil {
					adminAreaLevel2[i] = *val
				}
			}
			originFilter.AdminAreaLevel2 = adminAreaLevel2
		}
		if len(filter.Origin.AdminAreaLevel3) > 0 {
			adminAreaLevel3 := make([]string, len(filter.Origin.AdminAreaLevel3))
			for i, val := range filter.Origin.AdminAreaLevel3 {
				if val != nil {
					adminAreaLevel3[i] = *val
				}
			}
			originFilter.AdminAreaLevel3 = adminAreaLevel3
		}
		if len(filter.Origin.AdminAreaLevel4) > 0 {
			adminAreaLevel4 := make([]string, len(filter.Origin.AdminAreaLevel4))
			for i, val := range filter.Origin.AdminAreaLevel4 {
				if val != nil {
					adminAreaLevel4[i] = *val
				}
			}
			originFilter.AdminAreaLevel4 = adminAreaLevel4
		}

		if len(filter.Origin.ZipCodes) > 0 {
			zipCodes := make([]string, len(filter.Origin.ZipCodes))
			for i, zip := range filter.Origin.ZipCodes {
				if zip != nil {
					zipCodes[i] = *zip
				}
			}
			originFilter.ZipCodes = zipCodes
		}

		if filter.Origin.CoordinatesConfidence != nil {
			originFilter.CoordinatesConfidence = &domain.CoordinatesConfidenceLevelFilter{
				Min: filter.Origin.CoordinatesConfidence.Min,
				Max: filter.Origin.CoordinatesConfidence.Max,
			}
		}

		deliveryUnitsFilter.Origin = originFilter
	}

	// Filtros de Destination
	if filter.Destination != nil {
		destinationFilter := &domain.LocationFilter{}

		if len(filter.Destination.NodeReferences) > 0 {
			nodeReferences := make([]string, len(filter.Destination.NodeReferences))
			for i, ref := range filter.Destination.NodeReferences {
				if ref != nil {
					nodeReferences[i] = *ref
				}
			}
			destinationFilter.NodeReferences = nodeReferences
		}

		if len(filter.Destination.AddressLines) > 0 {
			addressLines := make([]string, len(filter.Destination.AddressLines))
			for i, line := range filter.Destination.AddressLines {
				if line != nil {
					addressLines[i] = *line
				}
			}
			destinationFilter.AddressLines = addressLines
		}

		if len(filter.Destination.AdminAreaLevel1) > 0 {
			adminAreaLevel1 := make([]string, len(filter.Destination.AdminAreaLevel1))
			for i, val := range filter.Destination.AdminAreaLevel1 {
				if val != nil {
					adminAreaLevel1[i] = *val
				}
			}
			destinationFilter.AdminAreaLevel1 = adminAreaLevel1
		}
		if len(filter.Destination.AdminAreaLevel2) > 0 {
			adminAreaLevel2 := make([]string, len(filter.Destination.AdminAreaLevel2))
			for i, val := range filter.Destination.AdminAreaLevel2 {
				if val != nil {
					adminAreaLevel2[i] = *val
				}
			}
			destinationFilter.AdminAreaLevel2 = adminAreaLevel2
		}
		if len(filter.Destination.AdminAreaLevel3) > 0 {
			adminAreaLevel3 := make([]string, len(filter.Destination.AdminAreaLevel3))
			for i, val := range filter.Destination.AdminAreaLevel3 {
				if val != nil {
					adminAreaLevel3[i] = *val
				}
			}
			destinationFilter.AdminAreaLevel3 = adminAreaLevel3
		}
		if len(filter.Destination.AdminAreaLevel4) > 0 {
			adminAreaLevel4 := make([]string, len(filter.Destination.AdminAreaLevel4))
			for i, val := range filter.Destination.AdminAreaLevel4 {
				if val != nil {
					adminAreaLevel4[i] = *val
				}
			}
			destinationFilter.AdminAreaLevel4 = adminAreaLevel4
		}

		if len(filter.Destination.ZipCodes) > 0 {
			zipCodes := make([]string, len(filter.Destination.ZipCodes))
			for i, zip := range filter.Destination.ZipCodes {
				if zip != nil {
					zipCodes[i] = *zip
				}
			}
			destinationFilter.ZipCodes = zipCodes
		}

		if filter.Destination.CoordinatesConfidence != nil {
			destinationFilter.CoordinatesConfidence = &domain.CoordinatesConfidenceLevelFilter{
				Min: filter.Destination.CoordinatesConfidence.Min,
				Max: filter.Destination.CoordinatesConfidence.Max,
			}
		}

		deliveryUnitsFilter.Destination = destinationFilter
	}

	// Filtros de PromisedDate
	if filter.PromisedDate != nil {
		promisedDateFilter := &domain.PromisedDateFilter{}

		if filter.PromisedDate.DateRange != nil {
			promisedDateFilter.DateRange = &domain.DateRangeFilter{
				StartDate: filter.PromisedDate.DateRange.StartDate,
				EndDate:   filter.PromisedDate.DateRange.EndDate,
			}
		}

		if filter.PromisedDate.TimeRange != nil {
			promisedDateFilter.TimeRange = &domain.TimeRangeFilter{
				StartTime: filter.PromisedDate.TimeRange.StartTime,
				EndTime:   filter.PromisedDate.TimeRange.EndTime,
			}
		}

		deliveryUnitsFilter.PromisedDate = promisedDateFilter
	}

	// Filtros de CollectAvailability
	if filter.CollectAvailability != nil {
		collectAvailabilityFilter := &domain.CollectAvailabilityFilter{}

		if len(filter.CollectAvailability.Dates) > 0 {
			dates := make([]string, len(filter.CollectAvailability.Dates))
			for i, date := range filter.CollectAvailability.Dates {
				if date != nil {
					dates[i] = *date
				}
			}
			collectAvailabilityFilter.Dates = dates
		}

		if filter.CollectAvailability.TimeRange != nil {
			collectAvailabilityFilter.TimeRange = &domain.TimeRangeFilter{
				StartTime: filter.CollectAvailability.TimeRange.StartTime,
				EndTime:   filter.CollectAvailability.TimeRange.EndTime,
			}
		}

		deliveryUnitsFilter.CollectAvailability = collectAvailabilityFilter
	}

	return deliveryUnitsFilter
}
