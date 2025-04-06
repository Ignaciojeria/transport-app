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
		orderStatus   domain.OrderStatus
		orderType     domain.OrderType
		originNode    domain.NodeInfo
		destNode      domain.NodeInfo
		originContact domain.Contact
		destContact   domain.Contact
		originAddress domain.AddressInfo
		destAddress   domain.AddressInfo
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Datos básicos para los tests
		orderStatus = domain.OrderStatus{
			Status: "pending",
		}

		orderType = domain.OrderType{
			Type:         "retail",
			Description:  "Entrega estándar",
			Organization: organization1,
		}

		originAddress = domain.AddressInfo{
			Organization: organization1,
			State:        "Región Metropolitana",
			AddressLine1: "Av. Providencia 1234",
			Location:     orb.Point{-70.6199, -33.4342},
		}

		destAddress = domain.AddressInfo{
			Organization: organization1,
			State:        "Región Metropolitana",
			AddressLine1: "Av. Las Condes 5678",
			Location:     orb.Point{-70.5714, -33.4012},
		}

		originContact = domain.Contact{
			FullName:     "Cliente Origen",
			PrimaryEmail: "origen@example.com",
			Organization: organization1,
		}

		destContact = domain.Contact{
			FullName:     "Cliente Destino",
			PrimaryEmail: "destino@example.com",
			Organization: organization1,
		}

		originNode = domain.NodeInfo{
			ReferenceID:  "node-origin-001",
			Name:         "Centro de Distribución Central",
			Organization: organization1,
			AddressInfo:  originAddress,
			Contact:      originContact,
		}

		destNode = domain.NodeInfo{
			ReferenceID:  "node-dest-001",
			Name:         "Punto de Entrega",
			Organization: organization1,
			AddressInfo:  destAddress,
			Contact:      destContact,
		}

		// Preparamos únicamente OrderStatus para los tests ya que es requerido
		err := connection.DB.WithContext(ctx).
			Table("order_statuses").
			Where("status = ?", orderStatus.Status).
			FirstOrCreate(&table.OrderStatus{
				Status: orderStatus.Status,
			}).Error
		Expect(err).ToNot(HaveOccurred())

		// Y también inicializamos el OrderStatus "in_progress" que usaremos más adelante
		err = connection.DB.WithContext(ctx).
			Table("order_statuses").
			Where("status = ?", "in_progress").
			FirstOrCreate(&table.OrderStatus{
				Status: "in_progress",
			}).Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert a new order if it doesn't exist", func() {
		order := domain.Order{
			Headers: domain.Headers{
				Organization: organization1,
				Commerce:     "Tienda Online",
				Consumer:     "Distribución Nacional",
			},
			ReferenceID:          "ORDER-001",
			OrderStatus:          orderStatus,
			OrderType:            orderType,
			Origin:               originNode,
			Destination:          destNode,
			DeliveryInstructions: "Dejar con el conserje",
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
		}

		upsert := NewUpsertOrder(connection)
		err := upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var dbOrder table.Order
		err = connection.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID()).
			First(&dbOrder).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbOrder.ReferenceID).To(Equal(string(order.ReferenceID)))
		Expect(dbOrder.DeliveryInstructions).To(Equal(order.DeliveryInstructions))
		Expect(dbOrder.OrganizationID).To(Equal(organization1.ID))
	})

	It("should update an existing order if delivery instructions changed", func() {
		order := domain.Order{
			Headers: domain.Headers{
				Organization: organization1,
				Commerce:     "Tienda Online",
				Consumer:     "Distribución Nacional",
			},
			ReferenceID:          "ORDER-002",
			OrderStatus:          orderStatus,
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
			Where("document_id = ?", order.DocID()).
			First(&dbOrder).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbOrder.DeliveryInstructions).To(Equal("Instrucciones modificadas"))
	})

	It("should update promised date fields if changed", func() {
		originalDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedDate := time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)

		order := domain.Order{
			Headers: domain.Headers{
				Organization: organization1,
				Commerce:     "Tienda Online",
				Consumer:     "Distribución Nacional",
			},
			ReferenceID: "ORDER-003",
			OrderStatus: orderStatus,
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
			Where("document_id = ?", order.DocID()).
			First(&dbOrder).Error
		Expect(err).ToNot(HaveOccurred())

		// Verificar que las fechas se actualizaron
		Expect(dbOrder.PromisedDateRangeStart.Format("2006-01-02")).To(Equal(updatedDate.Format("2006-01-02")))
		Expect(dbOrder.PromisedTimeRangeStart).To(Equal("14:00"))
		Expect(dbOrder.PromisedTimeRangeEnd).To(Equal("18:00"))
	})

	It("should update order status if changed", func() {
		// Nuevo estado para actualizar
		newStatus := domain.OrderStatus{
			Status: "in_progress",
		}

		order := domain.Order{
			Headers: domain.Headers{
				Organization: organization1,
				Commerce:     "Tienda Online",
				Consumer:     "Distribución Nacional",
			},
			ReferenceID: "ORDER-004",
			OrderStatus: orderStatus, // status pendiente
			OrderType:   orderType,
			Origin:      originNode,
			Destination: destNode,
		}

		upsert := NewUpsertOrder(connection)
		err := upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Actualizar el estado de la orden
		modifiedOrder := order
		modifiedOrder.OrderStatus = newStatus

		err = upsert(ctx, modifiedOrder)
		Expect(err).ToNot(HaveOccurred())

		var dbOrder table.Order
		err = connection.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID()).
			First(&dbOrder).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbOrder.OrderStatusDoc).To(Equal(newStatus.DocID().String()))
	})

	It("should update items if changed", func() {
		order := domain.Order{
			Headers: domain.Headers{
				Organization: organization1,
				Commerce:     "Tienda Online",
				Consumer:     "Distribución Nacional",
			},
			ReferenceID: "ORDER-005",
			OrderStatus: orderStatus,
			OrderType:   orderType,
			Origin:      originNode,
			Destination: destNode,
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
		}

		upsert := NewUpsertOrder(connection)
		err := upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Añadir un nuevo ítem
		modifiedOrder := order
		modifiedOrder.Items = []domain.Item{
			{
				Sku:         "SKU-001",
				Description: "Producto inicial",
				Quantity: domain.Quantity{
					QuantityNumber: 1,
					QuantityUnit:   "unidad",
				},
			},
			{
				Sku:         "SKU-002",
				Description: "Producto adicional",
				Quantity: domain.Quantity{
					QuantityNumber: 3,
					QuantityUnit:   "unidades",
				},
			},
		}

		err = upsert(ctx, modifiedOrder)
		Expect(err).ToNot(HaveOccurred())

		var dbOrder table.Order
		err = connection.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID()).
			First(&dbOrder).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(len(dbOrder.JSONItems)).To(Equal(2))
	})

	It("should fail if the orders table does not exist", func() {
		order := domain.Order{
			Headers: domain.Headers{
				Organization: organization1,
				Commerce:     "Tienda Online",
				Consumer:     "Distribución Nacional",
			},
			ReferenceID: "ORDER-ERROR",
			OrderStatus: orderStatus,
			OrderType:   orderType,
			Origin:      originNode,
			Destination: destNode,
		}

		upsert := NewUpsertOrder(noTablesContainerConnection)
		err := upsert(ctx, order)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("orders"))
	})

	// Este test se puede añadir al conjunto de tests de UpsertOrder
	It("should update document references when related entities change", func() {
		order := domain.Order{
			Headers: domain.Headers{
				Organization: organization1,
				Commerce:     "Tienda Online",
				Consumer:     "Distribución Nacional",
			},
			ReferenceID: "ORDER-DOC-TEST",
			OrderStatus: domain.OrderStatus{
				Status: "pending",
			},
			OrderType: domain.OrderType{
				Type:         "retail",
				Organization: organization1,
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
			Where("document_id = ?", order.DocID()).
			First(&initialOrder).Error
		Expect(err).ToNot(HaveOccurred())

		// Crear una nueva versión con todos los campos modificados
		modifiedOrder := order
		modifiedOrder.OrderStatus = domain.OrderStatus{
			Status: "in_progress",
		}
		modifiedOrder.OrderType = domain.OrderType{
			Type:         "express",
			Organization: organization1,
		}
		modifiedOrder.Headers = domain.Headers{
			Organization: organization1,
			Commerce:     "Nuevo Comercio",
			Consumer:     "Nuevo Consumidor",
		}
		modifiedOrder.Origin = domain.NodeInfo{
			ReferenceID:  "modified-origin-node-001",
			Name:         "Nuevo Nodo Origen",
			Organization: organization1,
			AddressInfo: domain.AddressInfo{
				Organization: organization1,
				AddressLine1: "Nueva dirección origen",
				Location:     orb.Point{-70.5, -33.5},
			},
			Contact: domain.Contact{
				FullName:     "Nuevo contacto origen",
				PrimaryEmail: "nuevo-origen@example.com",
				Organization: organization1,
			},
		}
		modifiedOrder.Destination = domain.NodeInfo{
			ReferenceID:  "modified-dest-node-002",
			Name:         "Nuevo Nodo Destino",
			Organization: organization1,
			AddressInfo: domain.AddressInfo{
				Organization: organization1,
				AddressLine1: "Nueva dirección destino",
				Location:     orb.Point{-70.1, -33.1},
			},
			Contact: domain.Contact{
				FullName:     "Nuevo contacto destino",
				PrimaryEmail: "nuevo-destino@example.com",
				Organization: organization1,
			},
		}

		err = upsert(ctx, modifiedOrder)
		Expect(err).ToNot(HaveOccurred())

		var updatedOrder table.Order
		err = connection.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID()).
			First(&updatedOrder).Error
		Expect(err).ToNot(HaveOccurred())

		// Validaciones de DocIDs actualizados
		Expect(updatedOrder.OrderStatusDoc).To(Equal(modifiedOrder.OrderStatus.DocID().String()))
		Expect(updatedOrder.OrderTypeDoc).To(Equal(modifiedOrder.OrderType.DocID().String()))
		Expect(updatedOrder.OrderHeadersDoc).To(Equal(modifiedOrder.Headers.DocID().String()))

		Expect(updatedOrder.OriginNodeInfoDoc).To(Equal(modifiedOrder.Origin.DocID().String()))
		Expect(updatedOrder.OriginContactDoc).To(Equal(modifiedOrder.Origin.Contact.DocID().String()))
		Expect(updatedOrder.OriginAddressInfoDoc).To(Equal(modifiedOrder.Origin.AddressInfo.DocID().String()))

		Expect(updatedOrder.DestinationNodeInfoDoc).To(Equal(modifiedOrder.Destination.DocID().String()))
		Expect(updatedOrder.DestinationContactDoc).To(Equal(modifiedOrder.Destination.Contact.DocID().String()))
		Expect(updatedOrder.DestinationAddressInfoDoc).To(Equal(modifiedOrder.Destination.AddressInfo.DocID().String()))
	})

})

func createMinimalOrder() domain.Order {
	return domain.Order{
		Headers: domain.Headers{
			Organization: organization1,
			Commerce:     "Tienda",
			Consumer:     "Cliente",
		},
		ReferenceID: "ORDER-FAIL-TEST",
		OrderStatus: domain.OrderStatus{Status: "pending"},
		OrderType:   domain.OrderType{Type: "retail", Organization: organization1},
	}
}
