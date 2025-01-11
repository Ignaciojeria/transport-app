package mapper

import (
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	"github.com/biter777/countries"
)

func MapOrderDomain(t table.Order) domain.Order {
	return domain.Order{
		ID:          t.ID,
		ReferenceID: domain.ReferenceID(t.ReferenceID),
		Organization: domain.Organization{
			ID:      t.OrganizationCountryID,
			Country: countries.ByName(t.OrganizationCountry.Country),
		},
		BusinessIdentifiers: domain.BusinessIdentifiers{
			Commerce: t.Commerce.Name,
			Consumer: t.Consumer.Name,
		},
		OrderStatus: domain.OrderStatus{
			ID:        t.OrderStatus.ID,
			Status:    t.OrderStatus.Status,
			CreatedAt: t.OrderStatus.CreatedAt.String(),
		},
		OrderType: domain.OrderType{
			Type:        t.OrderType.Type,
			Description: t.OrderType.Description,
		},
		References: mapOrderReferences(t.OrderReferences),
		Origin: domain.Origin{
			NodeInfo: domain.NodeInfo{
				ReferenceID: domain.ReferenceID(t.OriginNodeInfo.ReferenceID),
				Name:        t.OriginNodeInfo.Name,
				Type:        t.OriginNodeInfo.Type,
				Operator: domain.Operator{
					Contact: mapContactDomain(t.OriginNodeInfo.Operator.Contact),
					Type:    t.OriginNodeInfo.Operator.Type,
				},
				References: mapNodeReferences(t.OriginNodeInfo.NodeReferences),
			},
			AddressInfo: mapAddressInfo(t.OriginAddressInfo, t.OriginContact),
		},
		Destination: domain.Destination{
			DeliveryInstructions: t.DeliveryInstructions,
			NodeInfo: domain.NodeInfo{
				ReferenceID: domain.ReferenceID(t.DestinationNodeInfo.ReferenceID),
				Name:        t.DestinationNodeInfo.Name,
				Type:        t.DestinationNodeInfo.Type,
				Operator: domain.Operator{
					Contact: mapContactDomain(t.DestinationNodeInfo.Operator.Contact),
					Type:    t.DestinationNodeInfo.Operator.Type,
				},
				References: mapNodeReferences(t.DestinationNodeInfo.NodeReferences),
			},
			AddressInfo: mapAddressInfo(t.DestinationAddressInfo, t.DestinationContact),
		},
		Items:                   mapItems(t.JSONItems),
		Packages:                mapPackages(t.Packages),
		CollectAvailabilityDate: domain.CollectAvailabilityDate{Date: t.CollectAvailabilityDate},
		PromisedDate: domain.PromisedDate{
			DateRange: domain.DateRange{
				StartDate: t.PromisedDateRangeStart,
				EndDate:   t.PromisedDateRangeEnd,
			},
			TimeRange: domain.TimeRange{
				StartTime: t.PromisedTimeRangeStart,
				EndTime:   t.PromisedTimeRangeEnd,
			},
		},
		Visits:                mapVisits(t.Visits),
		TransportRequirements: mapTransportRequirements(t.TransportRequirements),
	}
}

func mapOrderReferences(refs []table.OrderReferences) []domain.Reference {
	result := make([]domain.Reference, len(refs))
	for i, ref := range refs {
		result[i] = domain.Reference{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return result
}

func mapContactDomain(t table.Contact) domain.Contact {
	return domain.Contact{
		FullName:   t.FullName,
		Email:      t.Email,
		Phone:      t.Phone,
		NationalID: t.NationalID,
		Documents:  mapDocumentsDomain(t.Documents),
	}
}

func mapDocumentsDomain(docs table.JSONDocuments) []domain.Document {
	result := make([]domain.Document, len(docs))
	for i, doc := range docs {
		result[i] = domain.Document{
			Value: doc.Value,
			Type:  doc.Type,
		}
	}
	return result
}

func mapAddressInfo(t table.AddressInfo, c table.Contact) domain.AddressInfo {
	return domain.AddressInfo{
		Contact:      mapContactDomain(c),
		State:        t.State,
		County:       t.County,
		Province:     t.Province,
		District:     t.District,
		AddressLine1: t.AddressLine1,
		AddressLine2: t.AddressLine2,
		AddressLine3: t.AddressLine3,
		Latitude:     t.Latitude,
		Longitude:    t.Longitude,
		ZipCode:      t.ZipCode,
		TimeZone:     t.TimeZone,
	}
}

func mapNodeReferences(refs []table.NodeReference) []domain.Reference {
	result := make([]domain.Reference, len(refs))
	for i, ref := range refs {
		result[i] = domain.Reference{
			Type:  ref.Type,
			Value: ref.Value,
		}
	}
	return result
}

func mapItems(items table.JSONItems) []domain.Item {
	result := make([]domain.Item, len(items))
	for i, item := range items {
		result[i] = domain.Item{
			ReferenceID:       domain.ReferenceID(item.ReferenceID),
			LogisticCondition: item.LogisticCondition,
			Quantity: domain.Quantity{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
			Insurance: domain.Insurance{
				UnitValue: item.Insurance.UnitValue,
				Currency:  item.Insurance.Currency,
			},
			Description: item.Description,
			Dimensions: domain.Dimensions{
				Height: item.Dimensions.Height,
				Width:  item.Dimensions.Width,
				Depth:  item.Dimensions.Depth,
				Unit:   item.Dimensions.Unit,
			},
			Weight: domain.Weight{
				Value: item.Weight.Value,
				Unit:  item.Weight.Unit,
			},
		}
	}
	return result
}

func mapPackages(packages []table.Package) []domain.Package {
	result := make([]domain.Package, len(packages))
	for i, pkg := range packages {
		result[i] = domain.Package{
			ID:  pkg.ID,
			Lpn: pkg.Lpn,
			Dimensions: domain.Dimensions{
				Height: pkg.Dimensions.Height,
				Width:  pkg.Dimensions.Width,
				Depth:  pkg.Dimensions.Depth,
				Unit:   pkg.Dimensions.Unit,
			},
			Weight: domain.Weight{
				Value: pkg.Weight.Value,
				Unit:  pkg.Weight.Unit,
			},
			Insurance: domain.Insurance{
				UnitValue: pkg.Insurance.UnitValue,
				Currency:  pkg.Insurance.Currency,
			},
			ItemReferences: mapItemReferences(pkg.JSONItems),
		}
	}
	return result
}

func mapItemReferences(items table.JSONItems) []domain.ItemReference {
	result := make([]domain.ItemReference, len(items))
	for i, item := range items {
		result[i] = domain.ItemReference{
			ReferenceID: domain.ReferenceID(item.ReferenceID),
			Quantity: domain.Quantity{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
		}
	}
	return result
}

func mapVisits(visits []table.Visit) []domain.Visit {
	result := make([]domain.Visit, len(visits))
	for i, visit := range visits {
		result[i] = domain.Visit{
			Date: visit.Date,
			TimeRange: domain.TimeRange{
				StartTime: visit.TimeRangeStart,
				EndTime:   visit.TimeRangeEnd,
			},
		}
	}
	return result
}

func mapTransportRequirements(req []byte) []domain.Reference {
	// Implementa según el formato del JSON almacenado en la base de datos
	return []domain.Reference{}
}
