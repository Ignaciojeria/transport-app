package request

import (
	"transport-app/app/domain"
)

type SearchOrdersResponse struct {
	UpsertOrderRequest  // Hereda todos los campos del request
	BusinessIdentifiers struct {
		Commerce string `json:"commerce"`
		Consumer string `json:"consumer"`
	} `json:"businessIdentifiers"`
	OrderStatus struct {
		ID        int64
		Status    string `json:"status"`
		CreatedAt string `json:"createdAt"`
	} `json:"orderStatus"`
}

func MapSearchOrdersResponse(orders []domain.Order) []SearchOrdersResponse {
	var responses []SearchOrdersResponse
	for _, order := range orders {
		response := SearchOrdersResponse{}
		response.
			withReferenceID(order.ReferenceID).
			withHeaders(order.Headers).
			withDeliveryInstructions(order.DeliveryInstructions).
			withOrderStatus(order.OrderStatus).
			withCollectAvailabilityDate(order.CollectAvailabilityDate).
			withOrigin(order.Origin).
			withDestination(order.Destination).
			withPromisedDate(order.PromisedDate).
			withPackages(order.Packages).
			withReferences(order.References).
			withTransportRequirements(order.TransportRequirements).
			withOrderType(order.OrderType).
			withItems(order.Items)
		responses = append(responses, response)
	}
	return responses
}

func (res *SearchOrdersResponse) withHeaders(headers domain.Headers) *SearchOrdersResponse {
	res.BusinessIdentifiers.Consumer = headers.Consumer
	res.BusinessIdentifiers.Commerce = headers.Commerce
	return res
}

func (res *SearchOrdersResponse) withReferenceID(referenceID domain.ReferenceID) *SearchOrdersResponse {
	res.ReferenceID = string(referenceID)
	return res
}

func (res *SearchOrdersResponse) withDeliveryInstructions(instructions string) *SearchOrdersResponse {
	res.Destination.DeliveryInstructions = string(instructions)
	return res
}

func (res *SearchOrdersResponse) withOrderStatus(orderStatus domain.OrderStatus) *SearchOrdersResponse {
	res.OrderStatus.ID = orderStatus.ID
	res.OrderStatus.Status = orderStatus.Status
	res.OrderStatus.CreatedAt = orderStatus.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	return res
}

func (res *SearchOrdersResponse) withOrigin(origin domain.NodeInfo) *SearchOrdersResponse {
	res.Origin.AddressInfo.AddressLine1 = origin.AddressInfo.AddressLine1
	res.Origin.AddressInfo.AddressLine2 = origin.AddressInfo.AddressLine2
	res.Origin.AddressInfo.AddressLine3 = origin.AddressInfo.AddressLine3
	res.Origin.AddressInfo.County = origin.AddressInfo.County
	res.Origin.AddressInfo.District = origin.AddressInfo.District
	res.Origin.AddressInfo.Province = origin.AddressInfo.Province
	res.Origin.AddressInfo.State = origin.AddressInfo.State
	res.Origin.AddressInfo.ZipCode = origin.AddressInfo.ZipCode
	res.Origin.AddressInfo.TimeZone = origin.AddressInfo.TimeZone
	res.Origin.AddressInfo.Latitude = origin.AddressInfo.Location[1]  // Latitud
	res.Origin.AddressInfo.Longitude = origin.AddressInfo.Location[0] // Longitud
	res.Origin.NodeInfo.ReferenceID = string(origin.ReferenceID)

	// ✅ Asegurar que se asigna el contacto de origen
	res.Origin.AddressInfo.Contact.FullName = origin.Contact.FullName
	res.Origin.AddressInfo.Contact.Email = origin.Contact.Email
	res.Origin.AddressInfo.Contact.Phone = origin.Contact.Phone
	res.Origin.AddressInfo.Contact.NationalID = origin.Contact.NationalID

	// ✅ Manejar documentos correctamente
	if len(origin.Contact.Documents) > 0 {
		res.Origin.AddressInfo.Contact.Documents = make([]struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}, len(origin.Contact.Documents))

		for i, doc := range origin.Contact.Documents {
			res.Origin.AddressInfo.Contact.Documents[i] = struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			}{
				Type:  doc.Type,
				Value: doc.Value,
			}
		}
	} else {
		res.Origin.AddressInfo.Contact.Documents = nil
	}

	return res
}

func (res *SearchOrdersResponse) withItems(items []domain.Item) *SearchOrdersResponse {
	res.UpsertOrderRequest.Items = make([]struct {
		Description string `json:"description"`
		Dimensions  struct {
			Depth  float64 `json:"depth"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		} `json:"insurance"`
		LogisticCondition string `json:"logisticCondition"`
		Quantity          struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		ReferenceID string `json:"referenceId"`
		Weight      struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"weight"`
	}, len(items))

	// Mapeo de cada item
	for i, item := range items {
		res.UpsertOrderRequest.Items[i] = struct {
			Description string `json:"description"`
			Dimensions  struct {
				Depth  float64 `json:"depth"`
				Height float64 `json:"height"`
				Unit   string  `json:"unit"`
				Width  float64 `json:"width"`
			} `json:"dimensions"`
			Insurance struct {
				Currency  string  `json:"currency"`
				UnitValue float64 `json:"unitValue"`
			} `json:"insurance"`
			LogisticCondition string `json:"logisticCondition"`
			Quantity          struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			ReferenceID string `json:"referenceId"`
			Weight      struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			} `json:"weight"`
		}{
			Description: item.Description,
			Dimensions: struct {
				Depth  float64 `json:"depth"`
				Height float64 `json:"height"`
				Unit   string  `json:"unit"`
				Width  float64 `json:"width"`
			}{
				Depth:  item.Dimensions.Depth,
				Height: item.Dimensions.Height,
				Unit:   item.Dimensions.Unit,
				Width:  item.Dimensions.Width,
			},
			Insurance: struct {
				Currency  string  `json:"currency"`
				UnitValue float64 `json:"unitValue"`
			}{
				Currency:  item.Insurance.Currency,
				UnitValue: item.Insurance.UnitValue,
			},
			LogisticCondition: item.LogisticCondition,
			Quantity: struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			}{
				QuantityNumber: item.Quantity.QuantityNumber,
				QuantityUnit:   item.Quantity.QuantityUnit,
			},
			ReferenceID: string(item.ReferenceID),
			Weight: struct {
				Unit  string  `json:"unit"`
				Value float64 `json:"value"`
			}{
				Unit:  item.Weight.Unit,
				Value: item.Weight.Value,
			},
		}
	}

	return res
}

func (res *SearchOrdersResponse) withDestination(destination domain.NodeInfo) *SearchOrdersResponse {
	res.Destination.AddressInfo.AddressLine1 = destination.AddressInfo.AddressLine1
	res.Destination.AddressInfo.AddressLine2 = destination.AddressInfo.AddressLine2
	res.Destination.AddressInfo.AddressLine3 = destination.AddressInfo.AddressLine3
	res.Destination.AddressInfo.County = destination.AddressInfo.County
	res.Destination.AddressInfo.District = destination.AddressInfo.District
	res.Destination.AddressInfo.Province = destination.AddressInfo.Province
	res.Destination.AddressInfo.State = destination.AddressInfo.State
	res.Destination.AddressInfo.ZipCode = destination.AddressInfo.ZipCode
	res.Destination.AddressInfo.TimeZone = destination.AddressInfo.TimeZone
	res.Destination.AddressInfo.Latitude = destination.AddressInfo.Location[1]  // Latitud
	res.Destination.AddressInfo.Longitude = destination.AddressInfo.Location[0] // Longitud
	res.Destination.NodeInfo.ReferenceID = string(destination.ReferenceID)

	// ✅ Asegurar que se asigna el contacto de destino
	res.Destination.AddressInfo.Contact.FullName = destination.Contact.FullName
	res.Destination.AddressInfo.Contact.Email = destination.Contact.Email
	res.Destination.AddressInfo.Contact.Phone = destination.Contact.Phone
	res.Destination.AddressInfo.Contact.NationalID = destination.Contact.NationalID

	// ✅ Manejar documentos correctamente
	if len(destination.Contact.Documents) > 0 {
		res.Destination.AddressInfo.Contact.Documents = make([]struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}, len(destination.Contact.Documents))

		for i, doc := range destination.Contact.Documents {
			res.Destination.AddressInfo.Contact.Documents[i] = struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			}{
				Type:  doc.Type,
				Value: doc.Value,
			}
		}
	} else {
		res.Destination.AddressInfo.Contact.Documents = nil // Para evitar JSON `null` en lugar de `[]`
	}

	return res
}

func (res *SearchOrdersResponse) withCollectAvailabilityDate(collectAvailabilityDate domain.CollectAvailabilityDate) *SearchOrdersResponse {
	// Convertir time.Time a string en formato yyyy-mm-dd
	res.CollectAvailabilityDate.Date = collectAvailabilityDate.Date.Format("2006-01-02")
	res.CollectAvailabilityDate.TimeRange.EndTime = collectAvailabilityDate.TimeRange.EndTime
	res.CollectAvailabilityDate.TimeRange.StartTime = collectAvailabilityDate.TimeRange.StartTime
	return res
}

func (res *SearchOrdersResponse) withPromisedDate(promisedDate domain.PromisedDate) *SearchOrdersResponse {
	// Convertir time.Time a string en formato yyyy-mm-dd
	res.PromisedDate.DateRange.StartDate = promisedDate.DateRange.StartDate.Format("2006-01-02")
	res.PromisedDate.DateRange.EndDate = promisedDate.DateRange.EndDate.Format("2006-01-02")
	res.PromisedDate.ServiceCategory = promisedDate.ServiceCategory
	res.PromisedDate.TimeRange.StartTime = promisedDate.TimeRange.StartTime
	res.PromisedDate.TimeRange.EndTime = promisedDate.TimeRange.EndTime
	return res
}

// Extrae la lógica de mapeo de referencias de ítems en una función separada
func mapItemReferences(itemReferences []domain.ItemReference) []struct {
	Quantity struct {
		QuantityNumber int    `json:"quantityNumber"`
		QuantityUnit   string `json:"quantityUnit"`
	} `json:"quantity"`
	ReferenceID string `json:"referenceId"`
} {
	mapped := make([]struct {
		Quantity struct {
			QuantityNumber int    `json:"quantityNumber"`
			QuantityUnit   string `json:"quantityUnit"`
		} `json:"quantity"`
		ReferenceID string `json:"referenceId"`
	}, len(itemReferences))

	for i, itemRef := range itemReferences {
		mapped[i] = struct {
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			ReferenceID string `json:"referenceId"`
		}{
			Quantity: struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			}{
				QuantityNumber: itemRef.Quantity.QuantityNumber,
				QuantityUnit:   itemRef.Quantity.QuantityUnit,
			},
			ReferenceID: string(itemRef.ReferenceID),
		}
	}

	return mapped
}

func (res *SearchOrdersResponse) withPackages(packages []domain.Package) *SearchOrdersResponse {
	// Reinicializa la lista de paquetes
	res.Packages = make([]struct {
		Dimensions struct {
			Depth  float64 `json:"depth"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
		} `json:"dimensions"`
		Insurance struct {
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		} `json:"insurance"`
		ItemReferences []struct {
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			ReferenceID string `json:"referenceId"`
		} `json:"itemReferences"`
		Lpn    string `json:"lpn"`
		Weight struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		} `json:"weight"`
	}, len(packages))

	// Mapeo de cada paquete
	for i, pkg := range packages {
		res.Packages[i].Dimensions = struct {
			Depth  float64 `json:"depth"`
			Height float64 `json:"height"`
			Unit   string  `json:"unit"`
			Width  float64 `json:"width"`
		}{
			Depth:  pkg.Dimensions.Depth,
			Height: pkg.Dimensions.Height,
			Unit:   pkg.Dimensions.Unit,
			Width:  pkg.Dimensions.Width,
		}

		res.Packages[i].Insurance = struct {
			Currency  string  `json:"currency"`
			UnitValue float64 `json:"unitValue"`
		}{
			Currency:  pkg.Insurance.Currency,
			UnitValue: pkg.Insurance.UnitValue,
		}

		res.Packages[i].ItemReferences = make([]struct {
			Quantity struct {
				QuantityNumber int    `json:"quantityNumber"`
				QuantityUnit   string `json:"quantityUnit"`
			} `json:"quantity"`
			ReferenceID string `json:"referenceId"`
		}, len(pkg.ItemReferences))

		for j, itemRef := range pkg.ItemReferences {
			res.Packages[i].ItemReferences[j] = struct {
				Quantity struct {
					QuantityNumber int    `json:"quantityNumber"`
					QuantityUnit   string `json:"quantityUnit"`
				} `json:"quantity"`
				ReferenceID string `json:"referenceId"`
			}{
				Quantity: struct {
					QuantityNumber int    `json:"quantityNumber"`
					QuantityUnit   string `json:"quantityUnit"`
				}{
					QuantityNumber: itemRef.Quantity.QuantityNumber,
					QuantityUnit:   itemRef.Quantity.QuantityUnit,
				},
				ReferenceID: string(itemRef.ReferenceID),
			}
		}

		res.Packages[i].Lpn = pkg.Lpn
		res.Packages[i].Weight = struct {
			Unit  string  `json:"unit"`
			Value float64 `json:"value"`
		}{
			Unit:  pkg.Weight.Unit,
			Value: pkg.Weight.Value,
		}
	}

	return res
}

func (res *SearchOrdersResponse) withReferences(references []domain.Reference) *SearchOrdersResponse {
	if res.References == nil {
		res.References = make([]struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}, 0)
	}

	res.References = res.References[:0]

	for _, ref := range references {
		refData := struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}{
			Type:  ref.Type,
			Value: ref.Value,
		}
		res.References = append(res.References, refData)
	}

	return res
}

func (res *SearchOrdersResponse) withTransportRequirements(requirements []domain.Reference) *SearchOrdersResponse {
	if res.TransportRequirements == nil {
		res.TransportRequirements = make([]struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}, 0)
	}

	res.TransportRequirements = res.TransportRequirements[:0]

	for _, req := range requirements {
		reqData := struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}{
			Type:  req.Type,
			Value: req.Value,
		}
		res.TransportRequirements = append(res.TransportRequirements, reqData)
	}

	return res
}

func (res *SearchOrdersResponse) withOrderType(orderType domain.OrderType) *SearchOrdersResponse {
	res.OrderType.Type = orderType.Type
	res.OrderType.Description = orderType.Description
	return res
}
