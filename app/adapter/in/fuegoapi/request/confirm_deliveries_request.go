package request

import (
	"context"
	"errors"
	"time"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"
)

type ConfirmDeliveriesRequest struct {
	ManualChange struct {
		PerformedBy string `json:"performedBy" example:"juan@example.com"`
		Reason      string `json:"reason" example:"Corrección tras reclamo de transporte"`
	} `json:"manualChange"`

	Carrier struct {
		Name       string `json:"name" example:"Transportes ABC"`
		NationalID string `json:"nationalID" example:"1234567890"`
	} `json:"carrier"`

	Driver struct {
		Email      string `json:"email" example:"juan@example.com"`
		NationalID string `json:"nationalID" example:"1234567890"`
	} `json:"driver"`

	Vehicle struct {
		Plate string `json:"plate" example:"ABC123"`
	} `json:"vehicle"`

	Route struct {
		ReferenceID    string `json:"referenceID" example:"ROUTE001"`
		SequenceNumber int    `json:"sequenceNumber" example:"1"`
	} `json:"route"`

	// Ahora cada DeliveryUnit tiene su propio estado de entrega
	DeliveryUnits []struct {
		// Identificadores de la orden padre
		OrderReferenceID string `json:"orderReferenceID" example:"ORD001"`

		BusinessIdentifiers struct {
			Commerce string `json:"commerce" example:"COMMERCE001"`
			Consumer string `json:"consumer" example:"CONSUMER001"`
		} `json:"businessIdentifiers"`

		Recipient struct {
			FullName   string `json:"fullName" example:"Juan Perez"`
			NationalID string `json:"nationalID" example:"1234567890"`
		} `json:"recipient"`

		// Cada unidad tiene su propio estado de entrega
		Delivery struct {
			HandledAt string `json:"handledAt" example:"2025-06-06T14:30:00Z"`
			Status    string `json:"status" example:"DELIVERED"` // DELIVERED, FAILED, PARTIAL
			Location  struct {
				Latitude  float64 `json:"latitude" example:"19.432607"`
				Longitude float64 `json:"longitude" example:"-99.133209"`
			} `json:"location"`
			Failure *struct {
				Detail      string `json:"detail" example:"no quiso recibir producto porque la caja estaba dañada"`
				Reason      string `json:"reason" example:"CLIENTE_RECHAZA_ENTREGA"`
				ReferenceID string `json:"referenceID" example:"1021"`
			} `json:"failure,omitempty"`
		} `json:"delivery"`

		// Fotos específicas para esta unidad de entrega
		EvidencePhotos []struct {
			TakenAt string `json:"takenAt" example:"2025-06-06T14:30:00Z"`
			Type    string `json:"type" example:"PACKAGE_PHOTO"`
			URL     string `json:"url" example:"https://ignaciojeria.github.io/"`
		} `json:"evidencePhotos"`

		// Información de la unidad de entrega
		Lpn   string `json:"lpn" example:"ABC123"`
		Items []struct {
			Sku         string `json:"sku" example:"SKU123"`
			Description string `json:"description" example:"Descripción del producto"`
			Quantity    int    `json:"quantity" example:"2"`
			// Para entregas parciales de items
			DeliveredQuantity int `json:"deliveredQuantity" example:"1"`
		} `json:"items"`
	} `json:"deliveryUnits"`
}

func (r ConfirmDeliveriesRequest) Map(ctx context.Context) domain.Route {

	route := domain.Route{
		ReferenceID: r.Route.ReferenceID,
		Vehicle: domain.Vehicle{
			Plate: r.Vehicle.Plate,
			Carrier: domain.Carrier{
				Name:       r.Carrier.Name,
				NationalID: r.Carrier.NationalID,
				Driver: domain.Driver{
					Email:      r.Driver.Email,
					NationalID: r.Driver.NationalID,
				},
			},
		},
	}

	// Crear un mapa para agrupar las unidades de entrega por orden
	ordersMap := make(map[string]*domain.Order)

	// Procesar cada unidad de entrega
	for _, du := range r.DeliveryUnits {
		orderRefID := du.OrderReferenceID

		// Crear o obtener la orden existente
		domainOrder, exists := ordersMap[orderRefID]
		if !exists {
			domainOrder = &domain.Order{
				Headers: domain.Headers{
					Consumer: du.BusinessIdentifiers.Consumer,
					Commerce: du.BusinessIdentifiers.Commerce,
					Channel:  sharedcontext.ChannelFromContext(ctx),
				},
				ReferenceID: domain.ReferenceID(du.OrderReferenceID),
				Destination: domain.NodeInfo{
					AddressInfo: domain.AddressInfo{
						Contact: domain.Contact{
							FullName:   du.Recipient.FullName,
							NationalID: du.Recipient.NationalID,
						},
					},
				},
				DeliveryUnits: make(domain.DeliveryUnits, 0),
			}
			ordersMap[orderRefID] = domainOrder
		}

		// Crear los items de la unidad de entrega
		items := make([]domain.Item, 0, len(du.Items))
		for _, item := range du.Items {
			items = append(items, domain.Item{
				Sku: item.Sku,
			})
		}

		// Parsear la fecha de manejo
		var handledAt time.Time
		if du.Delivery.HandledAt != "" {
			handledAt, _ = time.Parse(time.RFC3339, du.Delivery.HandledAt)
		}

		// Crear la unidad de entrega del dominio
		domainDeliveryUnit := domain.DeliveryUnit{
			Lpn:   du.Lpn,
			Items: items,
			ConfirmDelivery: domain.ConfirmDelivery{
				ManualChange: domain.ManualChange{
					PerformedBy: r.ManualChange.PerformedBy,
					Reason:      r.ManualChange.Reason,
				},
				HandledAt: handledAt,
				Latitude:  du.Delivery.Location.Latitude,
				Longitude: du.Delivery.Location.Longitude,
				EvidencePhotos: func() []domain.EvidencePhoto {
					photos := make([]domain.EvidencePhoto, 0, len(du.EvidencePhotos))
					for _, photo := range du.EvidencePhotos {
						takenAt, _ := time.Parse(time.RFC3339, photo.TakenAt)
						photos = append(photos, domain.EvidencePhoto{
							TakenAt: takenAt,
							Type:    photo.Type,
							URL:     photo.URL,
						})
					}
					return photos
				}(),
				Recipient: domain.Recipient{
					FullName:   du.Recipient.FullName,
					NationalID: du.Recipient.NationalID,
				},
			},
		}

		// Agregar información de fallo si existe
		if du.Delivery.Failure != nil {
			domainDeliveryUnit.ConfirmDelivery.NonDeliveryReason = domain.NonDeliveryReason{
				Reason:      du.Delivery.Failure.Reason,
				Details:     du.Delivery.Failure.Detail,
				ReferenceID: du.Delivery.Failure.ReferenceID,
			}
		}

		// Agregar la unidad de entrega a la orden
		domainOrder.DeliveryUnits = append(domainOrder.DeliveryUnits, domainDeliveryUnit)
	}

	// Convertir el mapa a slice de órdenes
	orders := make([]domain.Order, 0, len(ordersMap))
	for _, order := range ordersMap {
		orders = append(orders, *order)
	}

	route.Orders = orders

	return route
}

func (r ConfirmDeliveriesRequest) Validate() error {
	for _, du := range r.DeliveryUnits {
		// Si no hay LPN ni SKUs, retornar error
		if du.Lpn == "" && len(du.Items) == 0 {
			return errors.New("delivery unit must have either LPN or SKUs")
		}

		// Validar que cada unidad tenga un OrderReferenceID
		if du.OrderReferenceID == "" {
			return errors.New("delivery unit must have an orderReferenceID")
		}

		// Validar que tenga información de entrega
		if du.Delivery.HandledAt == "" {
			return errors.New("delivery unit must have a handledAt timestamp")
		}
	}
	return nil
}
