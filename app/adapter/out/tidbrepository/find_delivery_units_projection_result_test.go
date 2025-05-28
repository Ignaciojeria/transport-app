package tidbrepository

import (
	"context"
	"time"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"
	"transport-app/app/shared/projection/deliveryunits"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("FindDeliveryUnitsProjectionResult", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should return empty list when no delivery units exist", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			deliveryunits.NewProjection())
		results, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
	})

	It("should return delivery units when they exist", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred(), "Failed to create test tenant: %v", err)

		// Verify tenant was created
		var tenantCount int64
		err = conn.DB.WithContext(ctx).
			Table("tenants").
			Where("id = ?", tenant.ID).
			Count(&tenantCount).Error
		Expect(err).ToNot(HaveOccurred(), "Failed to verify tenant: %v", err)
		Expect(tenantCount).To(Equal(int64(1)), "Tenant was not created properly")

		destination := domain.AddressInfo{
			State:                "CA",
			Province:             "CA",
			District:             "CA",
			AddressLine1:         "123 Main St",
			AddressLine2:         "Apt 1",
			Location:             orb.Point{1, 1},
			TimeZone:             "America/Santiago",
			RequiresManualReview: true,
			CoordinateSource:     "geocoding",
			ZipCode:              "12345",
		}
		err = NewUpsertAddressInfo(conn)(ctx, destination)
		Expect(err).ToNot(HaveOccurred(), "Failed to upsert address info: %v", err)

		// Verify address was created
		var addressCount int64
		err = conn.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", destination.DocID(ctx)).
			Count(&addressCount).Error
		Expect(err).ToNot(HaveOccurred(), "Failed to verify address: %v", err)
		Expect(addressCount).To(Equal(int64(1)), "Address was not created properly.")

		fixedDate := time.Date(2025, 5, 26, 0, 0, 0, 0, time.UTC)
		order := domain.Order{
			ReferenceID:          "123",
			DeliveryInstructions: "Dejar en la puerta",
			Destination: domain.NodeInfo{
				AddressInfo: destination,
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{},
			},
			CollectAvailabilityDate: domain.CollectAvailabilityDate{
				Date: fixedDate,
				TimeRange: domain.TimeRange{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
			},
			PromisedDate: domain.PromisedDate{
				DateRange: domain.DateRange{
					StartDate: fixedDate,
					EndDate:   fixedDate.Add(24 * time.Hour),
				},
				TimeRange: domain.TimeRange{
					StartTime: "10:00",
					EndTime:   "17:00",
				},
				ServiceCategory: "STANDARD",
			},
		}
		err = NewUpsertOrder(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred(), "Failed to upsert order: %v", err)

		// Verify order was created
		var orderCount int64
		err = conn.DB.WithContext(ctx).
			Table("orders").
			Where("document_id = ?", order.DocID(ctx)).
			Count(&orderCount).Error
		Expect(err).ToNot(HaveOccurred(), "Failed to verify order: %v", err)
		Expect(orderCount).To(Equal(int64(1)), "Order was not created properly")

		// Create delivery units history
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: []domain.Order{
						order,
					},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred(), "Failed to upsert delivery units history: %v", err)

		// Verify delivery units history was created
		var historyCount int64
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_histories").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&historyCount).Error
		Expect(err).ToNot(HaveOccurred(), "Failed to verify delivery units history: %v", err)
		Expect(historyCount).To(Equal(int64(1)), "Delivery units history was not created properly")

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			deliveryunits.NewProjection())
		results, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                      "",
				projection.Channel().String():                          "",
				projection.CollectAvailabilityDate().String():          "",
				projection.CollectAvailabilityDateStartTime().String(): "",
				projection.CollectAvailabilityDateEndTime().String():   "",
				projection.PromisedDateDateRangeStartDate().String():   "",
				projection.PromisedDateDateRangeEndDate().String():     "",
				projection.PromisedDateTimeRangeStartTime().String():   "",
				projection.PromisedDateTimeRangeEndTime().String():     "",
				projection.DestinationAddressInfo().String():           "",
				projection.PromisedDateServiceCategory().String():      "",
				projection.DestinationAddressLine1().String():          "",
				projection.DestinationDistrict().String():              "",
				projection.DestinationLatitude().String():              "",
				projection.DestinationLongitude().String():             "",
				projection.DestinationProvince().String():              "",
				projection.DestinationState().String():                 "",
				projection.DestinationTimeZone().String():              "",
				projection.DestinationZipCode().String():               "",
				projection.DestinationAddressLine2().String():          "",
				projection.DeliveryInstructions().String():             "",
			},
		})
		Expect(err).ToNot(HaveOccurred(), "Failed to find delivery units: %v", err)
		Expect(results).To(HaveLen(1), "Expected 1 result, got %d", len(results))
		Expect(results[0].ID).To(BeNumerically(">", 0), "ID should be greater than 0")
		Expect(results[0].OrderReferenceID).To(Equal("123"), "Unexpected reference ID")
		Expect(results[0].OrderCollectAvailabilityDate).To(Equal("2025-05-26T00:00:00Z"), "Unexpected collect availability date")
		Expect(results[0].OrderCollectAvailabilityDateStartTime).To(Equal("09:00"), "Unexpected collect availability start time")
		Expect(results[0].OrderCollectAvailabilityDateEndTime).To(Equal("18:00"), "Unexpected collect availability end time")
		Expect(results[0].OrderPromisedDateStartDate).To(Equal("2025-05-26T00:00:00Z"), "Unexpected promised date start")
		Expect(results[0].OrderPromisedDateEndDate).To(Equal("2025-05-27T00:00:00Z"), "Unexpected promised date end")
		Expect(results[0].OrderPromisedDateStartTime).To(Equal("10:00"), "Unexpected promised time start")
		Expect(results[0].OrderPromisedDateEndTime).To(Equal("17:00"), "Unexpected promised time end")
		Expect(results[0].OrderPromisedDateServiceCategory).To(Equal("STANDARD"), "Unexpected service category")
		Expect(results[0].OrderDeliveryInstructions).To(Equal("Dejar en la puerta"), "Unexpected order delivery instructions")

		// Validaciones de Destination Address
		Expect(results[0].DestinationAddressLine1).To(Equal("123 Main St"), "Unexpected address line 1")
		Expect(results[0].DestinationAddressLine2).To(Equal("Apt 1"), "Unexpected address line 2")
		Expect(results[0].DestinationDistrict).To(Equal("CA"), "Unexpected district")
		Expect(results[0].DestinationLatitude).To(Equal(1.0), "Unexpected latitude")
		Expect(results[0].DestinationLongitude).To(Equal(1.0), "Unexpected longitude")
		Expect(results[0].DestinationProvince).To(Equal("CA"), "Unexpected province")
		Expect(results[0].DestinationState).To(Equal("CA"), "Unexpected state")
		Expect(results[0].DestinationTimeZone).To(Equal("America/Santiago"), "Unexpected timezone")
		Expect(results[0].DestinationZipCode).To(Equal("12345"), "Unexpected zip code")

	})

	It("should return delivery units with destination contact information", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear contacto
		contact := domain.Contact{
			FullName:     "John Doe",
			PrimaryEmail: "test@example.com",
			PrimaryPhone: "+56912345678",
			NationalID:   "12345678-9",
			AdditionalContactMethods: []domain.ContactMethod{
				{
					Type:  "whatsapp",
					Value: "+56987654321",
				},
				{
					Type:  "telegram",
					Value: "@johndoe",
				},
			},
			Documents: []domain.Document{
				{
					Type:  "dni",
					Value: "12345678-9",
				},
				{
					Type:  "passport",
					Value: "AB123456",
				},
			},
		}
		err = NewUpsertContact(conn)(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		destination := domain.AddressInfo{
			State:                "CA",
			Province:             "CA",
			District:             "CA",
			AddressLine1:         "123 Main St",
			Location:             orb.Point{1, 1},
			TimeZone:             "America/Santiago",
			RequiresManualReview: true,
			CoordinateSource:     "geocoding",
			ZipCode:              "12345",
			Contact:              contact,
		}
		err = NewUpsertAddressInfo(conn)(ctx, destination)
		Expect(err).ToNot(HaveOccurred())

		fixedDate := time.Date(2025, 5, 26, 0, 0, 0, 0, time.UTC)
		order := domain.Order{
			ReferenceID: "123",
			Destination: domain.NodeInfo{
				AddressInfo: destination,
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{},
			},
			CollectAvailabilityDate: domain.CollectAvailabilityDate{
				Date: fixedDate,
				TimeRange: domain.TimeRange{
					StartTime: "09:00",
					EndTime:   "18:00",
				},
			},
			PromisedDate: domain.PromisedDate{
				DateRange: domain.DateRange{
					StartDate: fixedDate,
					EndDate:   fixedDate.Add(24 * time.Hour),
				},
				TimeRange: domain.TimeRange{
					StartTime: "10:00",
					EndTime:   "17:00",
				},
				ServiceCategory: "STANDARD",
			},
		}
		err = NewUpsertOrder(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: []domain.Order{
						order,
					},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			deliveryunits.NewProjection())
		results, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.DestinationAddressInfo().String():              "",
				projection.DestinationContact().String():                  "",
				projection.DestinationContactEmail().String():             "",
				projection.DestinationContactFullName().String():          "",
				projection.DestinationContactNationalID().String():        "",
				projection.DestinationContactPhone().String():             "",
				projection.DestinationAdditionalContactMethods().String(): "",
				projection.DestinationContactDocuments().String():         "",
				projection.DestinationRequiresManualReview().String():     "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))

		// Validaciones de Destination Contact
		Expect(results[0].DestinationContactEmail).To(Equal("test@example.com"))
		Expect(results[0].DestinationContactFullName).To(Equal("John Doe"))
		Expect(results[0].DestinationContactNationalID).To(Equal("12345678-9"))
		Expect(results[0].DestinationContactPhone).To(Equal("+56912345678"))
		Expect(results[0].DestinationAdditionalContactMethods).To(Equal(table.JSONReference{
			{Type: "whatsapp", Value: "+56987654321"},
			{Type: "telegram", Value: "@johndoe"},
		}))
		Expect(results[0].DestinationContactDocuments).To(Equal(table.JSONReference{
			{Type: "dni", Value: "12345678-9"},
			{Type: "passport", Value: "AB123456"},
		}))
		Expect(results[0].DestinationRequiresManualReview).To(Equal(true))
	})

	It("should return delivery units with package information", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear delivery unit
		deliveryUnit := domain.DeliveryUnit{
			Lpn: "LPN123",
			Dimensions: domain.Dimensions{
				Length: 10,
				Width:  20,
				Height: 30,
				Unit:   "cm",
			},
			Weight: domain.Weight{
				Value: 5.5,
				Unit:  "kg",
			},
			Insurance: domain.Insurance{
				Currency:  "USD",
				UnitValue: 100.0,
			},
			Items: []domain.Item{
				{
					Sku:         "SKU123",
					Description: "Test Item",
					Quantity: domain.Quantity{
						QuantityNumber: 2,
						QuantityUnit:   "pcs",
					},
				},
			},
		}

		err = NewUpsertDeliveryUnits(conn)(ctx, []domain.DeliveryUnit{deliveryUnit})
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "123",
			Headers: domain.Headers{
				Commerce: "Test Commerce",
				Consumer: "Test Consumer",
			},
			DeliveryUnits: []domain.DeliveryUnit{
				deliveryUnit,
			},
		}
		err = NewUpsertOrder(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		err = NewUpsertOrderHeaders(conn)(ctx, order.Headers)
		Expect(err).ToNot(HaveOccurred())

		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: []domain.Order{
						order,
					},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			deliveryunits.NewProjection())
		results, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.DeliveryUnit().String():           "",
				projection.DeliveryUnitLPN().String():        "",
				projection.DeliveryUnitDimensions().String(): "",
				projection.DeliveryUnitWeight().String():     "",
				projection.DeliveryUnitInsurance().String():  "",
				projection.DeliveryUnitItems().String():      "",
				projection.Commerce().String():               "",
				projection.Consumer().String():               "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))

		// Validaciones de Delivery Unit
		Expect(results[0].LPN).To(Equal("LPN123"), "LPN incorrecto")
		Expect(results[0].JSONDimensions).To(Equal(table.JSONDimensions{
			Height: 30,
			Width:  20,
			Length: 10,
			Unit:   "cm",
		}), "Dimensiones incorrectas")
		Expect(results[0].JSONWeight).To(Equal(table.JSONWeight{
			WeightValue: 5.5,
			WeightUnit:  "kg",
		}), "Peso incorrecto")
		Expect(results[0].JSONInsurance).To(Equal(table.JSONInsurance{
			Currency:  "USD",
			UnitValue: 100.0,
		}), "Seguro incorrecto")

		// Validaciones de Items
		Expect(results[0].JSONItems).To(HaveLen(1), "Debería tener un item")
		Expect(results[0].JSONItems[0]).To(Equal(table.Items{
			Sku:            "SKU123",
			Description:    "Test Item",
			QuantityNumber: 2,
			QuantityUnit:   "pcs",
		}), "Item incorrecto")

		// Validaciones de Commerce y Consumer
		Expect(results[0].Commerce).To(Equal("Test Commerce"))
		Expect(results[0].Consumer).To(Equal("Test Consumer"))
	})

	It("should fail if database has no delivery_units_histories table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			noTablesContainerConnection,
			deliveryunits.NewProjection())
		_, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("delivery_units_histories"))
	})

	It("should return delivery units with order references", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear order references
		orderReferences := []domain.Reference{
			{
				Type:  "external",
				Value: "REF001",
			},
			{
				Type:  "internal",
				Value: "REF002",
			},
			{
				Type:  "tracking",
				Value: "REF003",
			},
		}

		order := domain.Order{
			ReferenceID: "123",
			References:  orderReferences,
			DeliveryUnits: []domain.DeliveryUnit{
				{},
			},
		}
		err = NewUpsertOrder(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		err = NewUpsertOrderReferences(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que las referencias se crearon
		var refCount int64
		err = conn.DB.WithContext(ctx).
			Table("order_references").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&refCount).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(refCount).To(Equal(int64(3)), "Order references were not created properly")

		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: []domain.Order{
						order,
					},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			deliveryunits.NewProjection())
		results, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.References().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))

		// Validar que las referencias se recuperaron correctamente
		Expect(results[0].OrderReferences).To(HaveLen(3))
		Expect(results[0].OrderReferences).To(ContainElements(
			table.Reference{Type: "external", Value: "REF001"},
			table.Reference{Type: "internal", Value: "REF002"},
			table.Reference{Type: "tracking", Value: "REF003"},
		))
	})

	It("should return delivery units with labels", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear delivery unit con etiquetas
		deliveryUnit := domain.DeliveryUnit{
			Lpn: "LPN123",
			Labels: []domain.Reference{
				{Type: "TYPE1", Value: "VALUE1"},
				{Type: "TYPE2", Value: "VALUE2"},
				{Type: "TYPE3", Value: "VALUE3"},
			},
		}

		order := domain.Order{
			ReferenceID: "123",
			DeliveryUnits: []domain.DeliveryUnit{
				deliveryUnit,
			},
		}

		err = NewUpsertDeliveryUnits(conn)(ctx, []domain.DeliveryUnit{deliveryUnit})
		Expect(err).ToNot(HaveOccurred())

		err = NewUpsertDeliveryUnitsLabels(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que las etiquetas se crearon correctamente
		var createdLabels []domain.Reference
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_labels").
			Where("delivery_unit_doc = ?", deliveryUnit.DocID(ctx)).
			Find(&createdLabels).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(createdLabels).To(HaveLen(3))

		err = NewUpsertOrder(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: []domain.Order{
						order,
					},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			deliveryunits.NewProjection())
		results, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.DeliveryUnitLabels().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))

		// Validar que las etiquetas se recuperaron correctamente
		Expect(results[0].DeliveryUnitLabels).To(HaveLen(3))
		Expect(results[0].DeliveryUnitLabels).To(ContainElements(
			table.Reference{Type: "TYPE1", Value: "VALUE1"},
			table.Reference{Type: "TYPE2", Value: "VALUE2"},
			table.Reference{Type: "TYPE3", Value: "VALUE3"},
		))
	})
})
