package mapper

import (
	"context"
	"time"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

func MapOrderToTable(ctx context.Context, order domain.Order) table.Order {
	//	orgCountryID := order.Organization.ID
	tbl := table.Order{
		/*
			DocumentID:     order.DocID().String(), // Agregar DocumentID
			ReferenceID:    string(order.ReferenceID),
			OrganizationID: order.Organization.ID,
		*/
		// Mapear IDs de documentos relacionados
		OrganizationID:         sharedcontext.TenantIDFromContext(ctx),
		ReferenceID:            string(order.ReferenceID),
		DocumentID:             order.DocID(ctx).String(),
		OrderHeadersDoc:        order.Headers.DocID(ctx).String(),
		OrderStatusDoc:         order.OrderStatus.DocID().String(),
		OrderTypeDoc:           order.OrderType.DocID(ctx).String(),
		OriginNodeInfoDoc:      order.Origin.DocID(ctx).String(),
		DestinationNodeInfoDoc: order.Destination.DocID(ctx).String(),
		ServiceCategory:        order.PromisedDate.ServiceCategory,
		// Si están disponibles, también mapear los contactos y direcciones
		OriginContactDoc:          order.Origin.AddressInfo.Contact.DocID(ctx).String(),
		DestinationContactDoc:     order.Destination.AddressInfo.Contact.DocID(ctx).String(),
		OriginAddressInfoDoc:      order.Origin.AddressInfo.DocID(ctx).String(),
		DestinationAddressInfoDoc: order.Destination.AddressInfo.DocID(ctx).String(),

		//OrderReferences:      mapReferencesToTable(order.References),
		DeliveryInstructions: order.DeliveryInstructions,

		CollectAvailabilityDate:           safePtrTime(order.CollectAvailabilityDate.Date),
		CollectAvailabilityTimeRangeStart: order.CollectAvailabilityDate.TimeRange.StartTime,
		CollectAvailabilityTimeRangeEnd:   order.CollectAvailabilityDate.TimeRange.EndTime,
		AddressLine2:                      order.AddressLine2,
		PromisedDateRangeStart:            safePtrTime(order.PromisedDate.DateRange.StartDate),
		PromisedDateRangeEnd:              safePtrTime(order.PromisedDate.DateRange.EndDate),
		PromisedTimeRangeStart:            order.PromisedDate.TimeRange.StartTime,
		PromisedTimeRangeEnd:              order.PromisedDate.TimeRange.EndTime,
		TransportRequirements:             mapTransportRequirementsToTable(order.TransportRequirements),
		//Packages:                          MapPackagesToTable(ctx, order.Packages),
	}
	return tbl
}

func safePtrTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil // Retorna nil si la fecha es vacía en el dominio
	}
	return &t
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

func MapPackagesToTable(ctx context.Context, packages []domain.Package) []table.Package {
	mapped := make([]table.Package, len(packages))
	for i, pkg := range packages {
		mapped[i] = table.Package{
			OrganizationID: sharedcontext.TenantIDFromContext(ctx),
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

func MapPackageToTable(ctx context.Context, pkg domain.Package, referenceId string) table.Package {
	return table.Package{
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),
		DocumentID:     pkg.DocID(ctx, referenceId).String(),
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
		JSONItems: mapItemsToTable(pkg.Items),
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

func MapAddressInfoToTable(ctx context.Context, address domain.AddressInfo) table.AddressInfo {
	return table.AddressInfo{
		OrganizationID: sharedcontext.TenantIDFromContext(ctx),

		State: address.State.String(),
		//	Locality:       address.Locality,
		District:     address.District.String(),
		AddressLine1: address.AddressLine1,
		//	AddressLine2:   address.AddressLine2,
		//	AddressLine3:   address.AddressLine3,
		//	RawAddress:     address.FullAddress(),
		DocumentID: string(address.DocID(ctx)),
		Latitude:   address.Location[1],
		Longitude:  address.Location[0],
		ZipCode:    address.ZipCode,
		Province:   address.Province.String(),
		TimeZone:   address.TimeZone,
	}
}
