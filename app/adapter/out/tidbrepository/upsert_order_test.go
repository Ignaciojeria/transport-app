package tidbrepository

import (
	"context"
	"time"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("UpsertOrder", func() {
	var (
		ctx           context.Context
		tenant        domain.Tenant
		orderType     domain.OrderType
		originNode    domain.NodeInfo
		destNode      domain.NodeInfo
		originAddress domain.AddressInfo
		destAddress   domain.AddressInfo
		loadStatuses  LoadStatuses
	)

	BeforeEach(func() {
		var err error
		// Create a new tenant for testing
		tenant, ctx, err = CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())

		// Initialize LoadStatuses
		loadStatuses = NewLoadStatuses(connection)
		err = loadStatuses()

		// Setup test data
		orderType = domain.OrderType{
			Type:        "retail",
			Description: "Test Order Type",
		}

		originAddress = domain.AddressInfo{
			AddressLine1: "Origin Address",
			Coordinates: domain.Coordinates{
				Point:  orb.Point{-70.6001, -33.4500},
				Source: "test",
				Confidence: domain.CoordinatesConfidence{
					Level:   1.0,
					Message: "Test confidence",
					Reason:  "Test data",
				},
			},
		}

		destAddress = domain.AddressInfo{
			AddressLine1: "Destination Address",
			Coordinates: domain.Coordinates{
				Point:  orb.Point{-70.6002, -33.4501},
				Source: "test",
				Confidence: domain.CoordinatesConfidence{
					Level:   1.0,
					Message: "Test confidence",
					Reason:  "Test data",
				},
			},
		}

		originNode = domain.NodeInfo{
			ReferenceID: "origin-001",
			Name:        "Origin Node",
			AddressInfo: originAddress,
		}

		destNode = domain.NodeInfo{
			ReferenceID: "dest-001",
			Name:        "Destination Node",
			AddressInfo: destAddress,
		}
	})

	It("should insert a new order if it doesn't exist", func() {
		order := domain.Order{
			Headers: domain.Headers{
				Commerce: "Tienda Online",
				Consumer: "Distribución Nacional",
			},
			ReferenceID:          "ORDER-001",
			OrderType:            orderType,
			Origin:               originNode,
			Destination:          destNode,
			DeliveryInstructions: "Dejar con el conserje",
			DeliveryUnits: []domain.DeliveryUnit{
				{
					Lpn: "PKG-001",
					Items: []domain.Item{
						{
							Sku:         "SKU-001",
							Description: "Producto de prueba",
							Quantity: domain.Quantity{
								QuantityNumber: 2,
								QuantityUnit:   "unidades",
							},
						},
					},
				},
			},
		}

		upsert := NewUpsertOrder(connection)
		err := upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var dbOrder table.Order
		err = connection.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID(ctx)).
			First(&dbOrder).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbOrder.ReferenceID).To(Equal(string(order.ReferenceID)))
		Expect(dbOrder.DeliveryInstructions).To(Equal(order.DeliveryInstructions))
		Expect(dbOrder.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update an existing order if delivery instructions changed", func() {
		order := domain.Order{
			Headers: domain.Headers{
				Commerce: "Tienda Online",
				Consumer: "Distribución Nacional",
			},
			ReferenceID:          "ORDER-002",
			OrderType:            orderType,
			Origin:               originNode,
			Destination:          destNode,
			DeliveryInstructions: "Instrucciones originales",
		}

		upsert := NewUpsertOrder(connection)
		err := upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Modificar las instrucciones de entrega
		modifiedOrder := order
		modifiedOrder.DeliveryInstructions = "Instrucciones modificadas"

		err = upsert(ctx, modifiedOrder)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que se haya actualizado
		var dbOrder table.Order
		err = connection.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID(ctx)).
			First(&dbOrder).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbOrder.DeliveryInstructions).To(Equal("Instrucciones modificadas"))
	})

	It("should update promised date fields if changed", func() {
		originalDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedDate := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

		order := domain.Order{
			Headers: domain.Headers{
				Commerce: "Tienda Online",
				Consumer: "Distribución Nacional",
			},
			ReferenceID: "ORDER-003",
			OrderType:   orderType,
			Origin:      originNode,
			Destination: destNode,
			PromisedDate: domain.PromisedDate{
				DateRange: domain.DateRange{
					StartDate: originalDate,
					EndDate:   originalDate,
				},
				TimeRange: domain.TimeRange{
					StartTime: "09:00",
					EndTime:   "12:00",
				},
			},
		}

		upsert := NewUpsertOrder(connection)
		err := upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Modificar las fechas prometidas
		modifiedOrder := order
		modifiedOrder.PromisedDate = domain.PromisedDate{
			DateRange: domain.DateRange{
				StartDate: updatedDate,
				EndDate:   updatedDate,
			},
			TimeRange: domain.TimeRange{
				StartTime: "14:00",
				EndTime:   "18:00",
			},
		}

		err = upsert(ctx, modifiedOrder)
		Expect(err).ToNot(HaveOccurred())

		var dbOrder table.Order
		err = connection.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID(ctx)).
			First(&dbOrder).Error
		Expect(err).ToNot(HaveOccurred())

		// Verificar que las fechas se actualizaron
		Expect(dbOrder.PromisedDateRangeStart.Format("2006-01-02")).To(Equal(updatedDate.Format("2006-01-02")))
		Expect(*dbOrder.PromisedTimeRangeStart).To(Equal("14:00:00"))
		Expect(*dbOrder.PromisedTimeRangeEnd).To(Equal("18:00:00"))
	})

	It("should update packages if changed", func() {
		order := domain.Order{
			Headers: domain.Headers{
				Commerce: "Tienda Online",
				Consumer: "Distribución Nacional",
			},
			ReferenceID: "ORDER-005",
			OrderType:   orderType,
			Origin:      originNode,
			Destination: destNode,
			DeliveryUnits: []domain.DeliveryUnit{
				{
					Lpn: "PKG-001",
					Items: []domain.Item{
						{
							Sku:         "SKU-001",
							Description: "Producto inicial",
							Quantity: domain.Quantity{
								QuantityNumber: 1,
								QuantityUnit:   "unidad",
							},
						},
					},
				},
			},
		}

		upsert := NewUpsertOrder(connection)
		err := upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Añadir un nuevo paquete
		modifiedOrder := order
		modifiedOrder.DeliveryUnits = []domain.DeliveryUnit{
			{
				Lpn: "PKG-001",
				Items: []domain.Item{
					{
						Sku:         "SKU-001",
						Description: "Producto inicial",
						Quantity: domain.Quantity{
							QuantityNumber: 1,
							QuantityUnit:   "unidad",
						},
					},
				},
			},
			{
				Lpn: "PKG-002",
				Items: []domain.Item{
					{
						Sku:         "SKU-002",
						Description: "Producto adicional",
						Quantity: domain.Quantity{
							QuantityNumber: 3,
							QuantityUnit:   "unidades",
						},
					},
				},
			},
		}

		err = upsert(ctx, modifiedOrder)
		Expect(err).ToNot(HaveOccurred())

		var dbOrder table.Order
		err = connection.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID(ctx)).
			First(&dbOrder).Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail if the orders table does not exist", func() {
		order := domain.Order{
			Headers: domain.Headers{
				Commerce: "Tienda Online",
				Consumer: "Distribución Nacional",
			},
			ReferenceID: "ORDER-ERROR",
			OrderType:   orderType,
			Origin:      originNode,
			Destination: destNode,
		}

		upsert := NewUpsertOrder(noTablesContainerConnection)
		err := upsert(ctx, order)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orders"))
	})

	It("should update document references when related entities change", func() {
		order := domain.Order{
			Headers: domain.Headers{
				Commerce: "Tienda Online",
				Consumer: "Distribución Nacional",
			},
			ReferenceID: "ORDER-DOC-TEST",
			OrderType: domain.OrderType{
				Type: "retail",
			},
			Origin:      originNode,
			Destination: destNode,
		}

		upsert := NewUpsertOrder(connection)
		err := upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var initialOrder table.Order
		err = connection.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID(ctx)).
			First(&initialOrder).Error
		Expect(err).ToNot(HaveOccurred())

		// Crear una nueva versión con todos los campos modificados
		modifiedOrder := order
		modifiedOrder.OrderType = domain.OrderType{
			Type: "express",
		}
		modifiedOrder.Headers = domain.Headers{
			Commerce: "Nuevo Comercio",
			Consumer: "Nuevo Consumidor",
		}
		modifiedOrder.Origin = domain.NodeInfo{
			ReferenceID: "modified-origin-node-001",
			Name:        "Nuevo Nodo Origen",
			AddressInfo: domain.AddressInfo{
				AddressLine1: "Nueva dirección origen",
				Coordinates: domain.Coordinates{
					Point:  orb.Point{-70.5, -33.5},
					Source: "test",
					Confidence: domain.CoordinatesConfidence{
						Level:   1.0,
						Message: "Test confidence",
						Reason:  "Test data",
					},
				},
			},
		}
		modifiedOrder.Destination = domain.NodeInfo{
			ReferenceID: "modified-dest-node-002",
			Name:        "Nuevo Nodo Destino",
			AddressInfo: domain.AddressInfo{
				AddressLine1: "Nueva dirección destino",
				Coordinates: domain.Coordinates{
					Point:  orb.Point{-70.1, -33.1},
					Source: "test",
					Confidence: domain.CoordinatesConfidence{
						Level:   1.0,
						Message: "Test confidence",
						Reason:  "Test data",
					},
				},
			},
		}

		err = upsert(ctx, modifiedOrder)
		Expect(err).ToNot(HaveOccurred())

		var updatedOrder table.Order
		err = connection.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID(ctx)).
			First(&updatedOrder).Error
		Expect(err).ToNot(HaveOccurred())

		Expect(updatedOrder.OrderTypeDoc).To(Equal(modifiedOrder.OrderType.DocID(ctx).String()))
		Expect(updatedOrder.OrderHeadersDoc).To(Equal(modifiedOrder.Headers.DocID(ctx).String()))

		Expect(updatedOrder.OriginNodeInfoDoc).To(Equal(modifiedOrder.Origin.DocID(ctx).String()))

		Expect(updatedOrder.OriginAddressInfoDoc).To(Equal(modifiedOrder.Origin.AddressInfo.DocID(ctx).String()))

		Expect(updatedOrder.DestinationNodeInfoDoc).To(Equal(modifiedOrder.Destination.DocID(ctx).String()))

		Expect(updatedOrder.DestinationAddressInfoDoc).To(Equal(modifiedOrder.Destination.AddressInfo.DocID(ctx).String()))
	})
})

func createMinimalOrder(ctx context.Context) domain.Order {
	return domain.Order{
		Headers: domain.Headers{
			Commerce: "Tienda",
			Consumer: "Cliente",
		},
		ReferenceID: "ORDER-FAIL-TEST",
		OrderType:   domain.OrderType{Type: "retail"},
	}
}
