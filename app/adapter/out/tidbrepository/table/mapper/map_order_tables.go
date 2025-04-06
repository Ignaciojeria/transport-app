package mapper

import (
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
)

func MapOrderToTable(order domain.Order) table.Order {
	orgCountryID := order.Organization.ID
		/*
	var planId *int64

	if order.Plan.ID != 0 { // Si el ID es distinto de 0, creamos un puntero
		planId = new(int64)
		*planId = order.Plan.ID
	}

	var routeID *int64
	if len(order.Plan.Routes) != 0 && order.Plan.Routes[0].ID != 0 { // Verifica que exista al menos una ruta y que su ID sea válido
		routeID = new(int64)
		*routeID = order.Plan.Routes[0].ID
	}*/
	tbl := table.Order{

		ReferenceID:    string(order.ReferenceID),
		OrganizationID: order.Organization.ID, // Completar según la lógica de negocio
		//CommerceID:            order.Commerce.ID,                        // Completar según la lógica de negocio
		//ConsumerID:            order.Consumer.ID,                        // Completar según la lógica de negocio
		OrderHeadersDoc: string(order.Headers.DocID()),

		//OrderType:       mapOrderTypeToTable(order.OrderType, orgCountryID),
		OrderReferences:      mapReferencesToTable(order.References),
		DeliveryInstructions: order.DeliveryInstructions,

		CollectAvailabilityDate:           safePtrTime(order.CollectAvailabilityDate.Date),
		CollectAvailabilityTimeRangeStart: order.CollectAvailabilityDate.TimeRange.StartTime,
		CollectAvailabilityTimeRangeEnd:   order.CollectAvailabilityDate.TimeRange.EndTime,
		PromisedDateRangeStart:            safePtrTime(order.PromisedDate.DateRange.StartDate),
		PromisedDateRangeEnd:              safePtrTime(order.PromisedDate.DateRange.EndDate),
		PromisedTimeRangeStart:            order.PromisedDate.TimeRange.StartTime,
		PromisedTimeRangeEnd:              order.PromisedDate.TimeRange.EndTime,
		JSONItems:                         mapItemsToTable(order.Items),
		//Visits:                            mapVisitsToTable(order.Visits),
		TransportRequirements: mapTransportRequirementsToTable(order.TransportRequirements),
		OrderHeaders:          mapHeadersToTable(order.Headers, orgCountryID),
		Packages:              MapPackagesToTable(order.Packages, orgCountryID),
	}
	return tbl
}

func safePtrTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil // Retorna nil si la fecha es vacía en el dominio
	}
	return &t
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
			Sku:               item.Sku,
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
				Length: item.Dimensions.Length,
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
			OrganizationID: orgCountryID,
			Lpn:            pkg.Lpn,
			JSONDimensions: table.JSONDimensions{
				Height: pkg.Dimensions.Height,
				Width:  pkg.Dimensions.Width,
				Length: pkg.Dimensions.Length,
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
			Sku:            item.Sku,
			QuantityNumber: item.Quantity.QuantityNumber,
			QuantityUnit:   item.Quantity.QuantityUnit,
		}
	}
	return mapped
}

func MapPackageToTable(pkg domain.Package) table.Package {
	return table.Package{
		OrganizationID: pkg.Organization.ID,
		DocumentID:     pkg.DocID().String(),
		Lpn:            pkg.Lpn,
		JSONDimensions: table.JSONDimensions{
			Height: pkg.Dimensions.Height,
			Width:  pkg.Dimensions.Width,
			Length: pkg.Dimensions.Length,
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

func mapHeadersToTable(c domain.Headers, orgCountryID int64) table.OrderHeaders {
	return table.OrderHeaders{
		Commerce: c.Commerce,
		Consumer: c.Consumer,
		Organization: table.Organization{
			ID: orgCountryID,
		},
	}
}

func mapOrderTypeToTable(t domain.OrderType, orgCountry int64) table.OrderType {
	return table.OrderType{
		OrganizationID: orgCountry,

		Type:        t.Type,
		Description: t.Description,
	}
}

func MapAddressInfoToTable(address domain.AddressInfo, orgCountry int64) table.AddressInfo {
	return table.AddressInfo{
		OrganizationID: orgCountry,

		State: address.State,
		//	Locality:       address.Locality,
		District:     address.District,
		AddressLine1: address.AddressLine1,
		//	AddressLine2:   address.AddressLine2,
		//	AddressLine3:   address.AddressLine3,
		//	RawAddress:     address.FullAddress(),
		DocumentID: string(address.DocID()),
		Latitude:   address.Location[1],
		Longitude:  address.Location[0],
		ZipCode:    address.ZipCode,
		Province:   address.Province,
		TimeZone:   address.TimeZone,
	}
}
