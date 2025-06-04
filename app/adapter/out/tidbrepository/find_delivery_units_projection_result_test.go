package tidbrepository

import (
	"context"
	"fmt"
	"time"
	"transport-app/app/adapter/out/tidbrepository/projectionresult"
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
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
		Expect(hasMore).To(BeFalse())
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
			State:        "CA",
			Province:     "CA",
			District:     "CA",
			AddressLine1: "123 Main St",
			AddressLine2: "Apt 1",
			Coordinates: domain.Coordinates{
				Point:  orb.Point{1, 1},
				Source: "geocoding",
				Confidence: domain.CoordinatesConfidence{
					Level:   0.8,
					Message: "Medium confidence",
					Reason:  "Geocoding service",
				},
			},
			TimeZone: "America/Santiago",
			ZipCode:  "12345",
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
			Table("delivery_units_status_histories").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&historyCount).Error
		Expect(err).ToNot(HaveOccurred(), "Failed to verify delivery units history: %v", err)
		Expect(historyCount).To(Equal(int64(1)), "Delivery units history was not created properly")

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			deliveryunits.NewProjection())
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                             "",
				projection.Channel().String():                                 "",
				projection.CollectAvailabilityDate().String():                 "",
				projection.CollectAvailabilityDateStartTime().String():        "",
				projection.CollectAvailabilityDateEndTime().String():          "",
				projection.PromisedDateDateRangeStartDate().String():          "",
				projection.PromisedDateDateRangeEndDate().String():            "",
				projection.PromisedDateTimeRangeStartTime().String():          "",
				projection.PromisedDateTimeRangeEndTime().String():            "",
				projection.DestinationAddressInfo().String():                  "",
				projection.PromisedDateServiceCategory().String():             "",
				projection.DestinationAddressLine1().String():                 "",
				projection.DestinationDistrict().String():                     "",
				projection.DestinationCoordinatesLatitude().String():          "",
				projection.DestinationCoordinatesLongitude().String():         "",
				projection.DestinationCoordinatesSource().String():            "",
				projection.DestinationCoordinatesConfidenceLevel().String():   "",
				projection.DestinationCoordinatesConfidenceMessage().String(): "",
				projection.DestinationCoordinatesConfidenceReason().String():  "",
				projection.DestinationProvince().String():                     "",
				projection.DestinationState().String():                        "",
				projection.DestinationTimeZone().String():                     "",
				projection.DestinationZipCode().String():                      "",
				projection.DestinationAddressLine2().String():                 "",
				projection.DeliveryInstructions().String():                    "",
			},
		})
		Expect(err).ToNot(HaveOccurred(), "Failed to find delivery units: %v", err)
		Expect(results).To(HaveLen(1), "Expected 1 result, got %d", len(results))
		Expect(hasMore).To(BeFalse())
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
		Expect(results[0].DestinationCoordinatesLatitude).To(Equal(1.0), "Unexpected latitude")
		Expect(results[0].DestinationCoordinatesLongitude).To(Equal(1.0), "Unexpected longitude")
		Expect(results[0].DestinationCoordinatesSource).To(Equal("geocoding"), "Unexpected source")
		Expect(results[0].DestinationCoordinatesConfidenceLevel).To(Equal(0.8), "Unexpected confidence level")
		Expect(results[0].DestinationCoordinatesConfidenceMessage).To(Equal("Medium confidence"), "Unexpected confidence message")
		Expect(results[0].DestinationCoordinatesConfidenceReason).To(Equal("Geocoding service"), "Unexpected confidence reason")
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
			State:        "CA",
			Province:     "CA",
			District:     "CA",
			AddressLine1: "123 Main St",
			Coordinates: domain.Coordinates{
				Point:  orb.Point{1, 1},
				Source: "geocoding",
				Confidence: domain.CoordinatesConfidence{
					Level:   0.8,
					Message: "Medium confidence",
					Reason:  "Geocoding service",
				},
			},
			TimeZone: "America/Santiago",
			ZipCode:  "12345",
			Contact:  contact,
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
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.DestinationAddressInfo().String():              "",
				projection.DestinationContact().String():                  "",
				projection.DestinationContactEmail().String():             "",
				projection.DestinationContactFullName().String():          "",
				projection.DestinationContactNationalID().String():        "",
				projection.DestinationContactPhone().String():             "",
				projection.DestinationAdditionalContactMethods().String(): "",
				projection.DestinationContactDocuments().String():         "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())

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
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
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
		Expect(hasMore).To(BeFalse())

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

	It("should return delivery units with extra fields", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear order con extra fields
		order := domain.Order{
			ReferenceID: "123",
			ExtraFields: map[string]string{
				"priority": "high",
				"notes":    "Handle with care",
				"tags":     "fragile,express",
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{},
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
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ExtraFields().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())

		// Validar que los extra fields se recuperaron correctamente
		Expect(results[0].ExtraFields).To(Equal(table.JSONMap{
			"priority": "high",
			"notes":    "Handle with care",
			"tags":     "fragile,express",
		}))
	})

	It("should return delivery units with order type", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear order type
		orderType := domain.OrderType{
			Type:        "EXPRESS",
			Description: "Entrega express",
		}
		err = NewUpsertOrderType(conn)(ctx, orderType)
		Expect(err).ToNot(HaveOccurred())

		// Crear order con order type
		order := domain.Order{
			ReferenceID: "123",
			OrderType:   orderType,
			DeliveryUnits: []domain.DeliveryUnit{
				{},
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
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.OrderType().String():            "",
				projection.OrderTypeType().String():        "",
				projection.OrderTypeDescription().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())

		// Validar que el order type se recuperó correctamente
		Expect(results[0].OrderType).To(Equal("EXPRESS"))
		Expect(results[0].OrderTypeDescription).To(Equal("Entrega express"))
	})

	It("should return delivery units with status", func() {
		// Create a new tenant for this test
		err := NewLoadStatuses(conn)()
		Expect(err).ToNot(HaveOccurred())

		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear order
		order := domain.Order{
			ReferenceID: "123",
			DeliveryUnits: []domain.DeliveryUnit{
				{
					Status: domain.Status{
						Status: domain.StatusPlanned,
					},
				},
			},
		}
		err = NewUpsertOrder(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Crear delivery units history con status
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
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.Status().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())

		// Validar que el status se recuperó correctamente
		Expect(results[0].Status).To(Equal(domain.StatusPlanned))
	})

	It("should fail if database has no delivery_units_status_histories table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			noTablesContainerConnection,
			deliveryunits.NewProjection())
		_, _, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("delivery_units_status_histories"))
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
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.References().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())

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
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.DeliveryUnitLabels().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())

		// Validar que las etiquetas se recuperaron correctamente
		Expect(results[0].DeliveryUnitLabels).To(HaveLen(3))
		Expect(results[0].DeliveryUnitLabels).To(ContainElements(
			table.Reference{Type: "TYPE1", Value: "VALUE1"},
			table.Reference{Type: "TYPE2", Value: "VALUE2"},
			table.Reference{Type: "TYPE3", Value: "VALUE3"},
		))
	})

	It("should return delivery units with origin information", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear contacto de origen
		originContact := domain.Contact{
			FullName:     "Jane Smith",
			PrimaryEmail: "origin@example.com",
			PrimaryPhone: "+56987654321",
			NationalID:   "87654321-9",
			AdditionalContactMethods: []domain.ContactMethod{
				{
					Type:  "whatsapp",
					Value: "+56912345678",
				},
				{
					Type:  "telegram",
					Value: "@janesmith",
				},
			},
			Documents: []domain.Document{
				{
					Type:  "dni",
					Value: "87654321-9",
				},
				{
					Type:  "passport",
					Value: "CD876543",
				},
			},
		}
		err = NewUpsertContact(conn)(ctx, originContact)
		Expect(err).ToNot(HaveOccurred())

		// Crear dirección de origen
		origin := domain.AddressInfo{
			State:        "NY",
			Province:     "NY",
			District:     "Manhattan",
			AddressLine1: "456 Park Ave",
			AddressLine2: "Suite 100",
			Coordinates: domain.Coordinates{
				Point:  orb.Point{2, 2},
				Source: "geocoding",
				Confidence: domain.CoordinatesConfidence{
					Level:   0.9,
					Message: "High confidence",
					Reason:  "Geocoding service",
				},
			},
			TimeZone: "America/New_York",
			ZipCode:  "10022",
			Contact:  originContact,
		}
		err = NewUpsertAddressInfo(conn)(ctx, origin)
		Expect(err).ToNot(HaveOccurred())

		// Crear order con información de origen
		order := domain.Order{
			ReferenceID: "123",
			Origin: domain.NodeInfo{
				AddressInfo: origin,
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{},
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
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.OriginAddressInfo().String():                  "",
				projection.OriginAddressLine1().String():                 "",
				projection.OriginAddressLine2().String():                 "",
				projection.OriginDistrict().String():                     "",
				projection.OriginProvince().String():                     "",
				projection.OriginState().String():                        "",
				projection.OriginZipCode().String():                      "",
				projection.OriginCoordinatesLatitude().String():          "",
				projection.OriginCoordinatesLongitude().String():         "",
				projection.OriginCoordinatesSource().String():            "",
				projection.OriginCoordinatesConfidenceLevel().String():   "",
				projection.OriginCoordinatesConfidenceMessage().String(): "",
				projection.OriginCoordinatesConfidenceReason().String():  "",
				projection.OriginTimeZone().String():                     "",
				projection.OriginContact().String():                      "",
				projection.OriginContactEmail().String():                 "",
				projection.OriginContactFullName().String():              "",
				projection.OriginContactNationalID().String():            "",
				projection.OriginContactPhone().String():                 "",
				projection.OriginContactMethods().String():               "",
				projection.OriginDocuments().String():                    "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())

		// Validaciones de Origin Address
		Expect(results[0].OriginAddressLine1).To(Equal("456 Park Ave"))
		Expect(results[0].OriginAddressLine2).To(Equal("Suite 100"))
		Expect(results[0].OriginDistrict).To(Equal("Manhattan"))
		Expect(results[0].OriginProvince).To(Equal("NY"))
		Expect(results[0].OriginState).To(Equal("NY"))
		Expect(results[0].OriginZipCode).To(Equal("10022"))
		Expect(results[0].OriginTimeZone).To(Equal("America/New_York"))

		// Validaciones de Origin Coordinates
		Expect(results[0].OriginCoordinatesLatitude).To(Equal(2.0))
		Expect(results[0].OriginCoordinatesLongitude).To(Equal(2.0))
		Expect(results[0].OriginCoordinatesSource).To(Equal("geocoding"))
		Expect(results[0].OriginCoordinatesConfidenceLevel).To(Equal(0.9))
		Expect(results[0].OriginCoordinatesConfidenceMessage).To(Equal("High confidence"))
		Expect(results[0].OriginCoordinatesConfidenceReason).To(Equal("Geocoding service"))

		// Validaciones de Origin Contact
		Expect(results[0].OriginContactEmail).To(Equal("origin@example.com"))
		Expect(results[0].OriginContactFullName).To(Equal("Jane Smith"))
		Expect(results[0].OriginContactNationalID).To(Equal("87654321-9"))
		Expect(results[0].OriginContactPhone).To(Equal("+56987654321"))
		Expect(results[0].OriginAdditionalContactMethods).To(Equal(table.JSONReference{
			{Type: "whatsapp", Value: "+56912345678"},
			{Type: "telegram", Value: "@janesmith"},
		}))
		Expect(results[0].OriginContactDocuments).To(Equal(table.JSONReference{
			{Type: "dni", Value: "87654321-9"},
			{Type: "passport", Value: "CD876543"},
		}))
	})

	It("should order results correctly based on pagination direction", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear varios orders con diferentes IDs para probar el ordenamiento
		orders := []domain.Order{
			{
				ReferenceID: "order1",
				DeliveryUnits: []domain.DeliveryUnit{
					{},
				},
			},
			{
				ReferenceID: "order2",
				DeliveryUnits: []domain.DeliveryUnit{
					{},
				},
			},
			{
				ReferenceID: "order3",
				DeliveryUnits: []domain.DeliveryUnit{
					{},
				},
			},
		}

		// Insertar los orders
		for _, order := range orders {
			err = NewUpsertOrder(conn)(ctx, order)
			Expect(err).ToNot(HaveOccurred())
		}

		// Crear delivery units history
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: orders,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			deliveryunits.NewProjection())

		// Probar paginación hacia adelante
		forwardResults, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			Pagination: domain.Pagination{
				First: ptr(1),
			},
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(forwardResults).To(HaveLen(1))
		Expect(hasMore).To(BeTrue(), "Debería indicar que hay más resultados")
		// Verificar que tenemos el primer order
		Expect(forwardResults[0].OrderReferenceID).To(Equal("order1"))

		// Probar paginación hacia atrás
		backwardResults, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			Pagination: domain.Pagination{
				Last: ptr(1),
			},
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(backwardResults).To(HaveLen(1))
		Expect(hasMore).To(BeTrue(), "Debería indicar que hay más resultados")
		// Verificar que tenemos el último order (el repositorio aplica Reversed() internamente)
		Expect(backwardResults[0].OrderReferenceID).To(Equal("order3"))

		// Probar sin paginación - debería devolver todos los resultados en orden ascendente
		noPaginationResults, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(noPaginationResults).To(HaveLen(3))
		Expect(hasMore).To(BeFalse(), "No debería indicar que hay más resultados cuando se obtienen todos")
		// Verificar que tenemos todos los orders en orden ascendente
		Expect(noPaginationResults[0].OrderReferenceID).To(Equal("order1"))
		Expect(noPaginationResults[1].OrderReferenceID).To(Equal("order2"))
		Expect(noPaginationResults[2].OrderReferenceID).To(Equal("order3"))

		// Probar paginación con límite igual al total de resultados
		exactLimitResults, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			Pagination: domain.Pagination{
				First: ptr(3),
			},
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(exactLimitResults).To(HaveLen(3))
		Expect(hasMore).To(BeFalse(), "No debería indicar que hay más resultados cuando el límite es igual al total")
	})

	It("should filter delivery units by references", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Create an order with references
		order := domain.Order{
			ReferenceID: "REF-001",
			References: []domain.Reference{
				{Type: "TRACKING", Value: "TRK-123"},
				{Type: "EXTERNAL", Value: "EXT-456"},
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{},
			},
		}
		err = NewUpsertOrder(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Create another order with different references
		order2 := domain.Order{
			ReferenceID: "REF-002",
			References: []domain.Reference{
				{Type: "TRACKING", Value: "TRK-789"},
				{Type: "EXTERNAL", Value: "EXT-012"},
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{},
			},
		}
		err = NewUpsertOrder(conn)(ctx, order2)
		Expect(err).ToNot(HaveOccurred())

		err = NewUpsertOrderReferences(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())
		err = NewUpsertOrderReferences(conn)(ctx, order2)
		Expect(err).ToNot(HaveOccurred())

		// Create delivery units history for both orders
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: []domain.Order{order, order2},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		// Test filtering by a single reference
		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			projection)

		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
				projection.References().String():  "",
			},
			Order: &domain.OrderFilter{
				References: []domain.ReferenceFilter{
					{Type: "TRACKING", Value: "TRK-123"},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())
		Expect(results[0].OrderReferenceID).To(Equal("REF-001"))

		// Test filtering by multiple references
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
				projection.References().String():  "",
			},
			Order: &domain.OrderFilter{
				References: []domain.ReferenceFilter{
					{Type: "EXTERNAL", Value: "EXT-012"},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())
		Expect(results[0].OrderReferenceID).To(Equal("REF-002"))

		// Test filtering with non-existent reference
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
				projection.References().String():  "",
			},
			Order: &domain.OrderFilter{
				References: []domain.ReferenceFilter{
					{Type: "TRACKING", Value: "NON-EXISTENT"},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
		Expect(hasMore).To(BeFalse())
	})

	It("should filter delivery units by reference ids", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Create orders with different reference IDs
		orders := []domain.Order{
			{
				ReferenceID: "REF-001",
				DeliveryUnits: []domain.DeliveryUnit{
					{},
				},
			},
			{
				ReferenceID: "REF-002",
				DeliveryUnits: []domain.DeliveryUnit{
					{},
				},
			},
			{
				ReferenceID: "REF-003",
				DeliveryUnits: []domain.DeliveryUnit{
					{},
				},
			},
		}

		// Insert orders
		for _, order := range orders {
			err = NewUpsertOrder(conn)(ctx, order)
			Expect(err).ToNot(HaveOccurred())
		}

		// Create delivery units history
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: orders,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			projection)

		// Test filtering by a single reference ID
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
			Order: &domain.OrderFilter{
				ReferenceIds: []string{"REF-001"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())
		Expect(results[0].OrderReferenceID).To(Equal("REF-001"))

		// Test filtering by multiple reference IDs
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
			Order: &domain.OrderFilter{
				ReferenceIds: []string{"REF-001", "REF-003"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(2))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-001")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-003")),
		))

		// Test filtering with non-existent reference ID
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
			Order: &domain.OrderFilter{
				ReferenceIds: []string{"NON-EXISTENT"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
		Expect(hasMore).To(BeFalse())
	})

	It("should filter delivery units by labels", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear delivery units con diferentes etiquetas
		deliveryUnit1 := domain.DeliveryUnit{
			Lpn: "LPN001",
			Labels: []domain.Reference{
				{Type: "PRIORITY", Value: "HIGH"},
				{Type: "FRAGILE", Value: "YES"},
			},
		}

		deliveryUnit2 := domain.DeliveryUnit{
			Lpn: "LPN002",
			Labels: []domain.Reference{
				{Type: "PRIORITY", Value: "LOW"},
				{Type: "TEMPERATURE", Value: "COLD"},
			},
		}

		deliveryUnit3 := domain.DeliveryUnit{
			Lpn: "LPN003",
			Labels: []domain.Reference{
				{Type: "PRIORITY", Value: "HIGH"},
				{Type: "TEMPERATURE", Value: "HOT"},
			},
		}

		// Crear orders con las delivery units
		orders := []domain.Order{
			{
				ReferenceID: "REF-001",
				DeliveryUnits: []domain.DeliveryUnit{
					deliveryUnit1,
				},
			},
			{
				ReferenceID: "REF-002",
				DeliveryUnits: []domain.DeliveryUnit{
					deliveryUnit2,
				},
			},
			{
				ReferenceID: "REF-003",
				DeliveryUnits: []domain.DeliveryUnit{
					deliveryUnit3,
				},
			},
		}

		// Insertar delivery units
		err = NewUpsertDeliveryUnits(conn)(ctx, []domain.DeliveryUnit{deliveryUnit1, deliveryUnit2, deliveryUnit3})
		Expect(err).ToNot(HaveOccurred())

		// Insertar orders
		for _, order := range orders {
			err = NewUpsertOrder(conn)(ctx, order)
			Expect(err).ToNot(HaveOccurred())
		}

		// Insertar labels
		for _, order := range orders {
			err = NewUpsertDeliveryUnitsLabels(conn)(ctx, order)
			Expect(err).ToNot(HaveOccurred())
		}

		// Crear delivery units history
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: orders,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			projection)

		// Test filtrando por una sola etiqueta
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():        "",
				projection.DeliveryUnitLabels().String(): "",
			},
			DeliveryUnit: &domain.DeliveryUnitFilter{
				Labels: []domain.LabelFilter{
					{Type: "PRIORITY", Value: "HIGH"},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(2))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-001")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-003")),
		))

		// Test filtrando por múltiples etiquetas
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():        "",
				projection.DeliveryUnitLabels().String(): "",
			},
			DeliveryUnit: &domain.DeliveryUnitFilter{
				Labels: []domain.LabelFilter{
					{Type: "PRIORITY", Value: "HIGH"},
					{Type: "TEMPERATURE", Value: "HOT"},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(2))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-001")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-003")),
		))

		// Test filtrando por etiqueta no existente
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():        "",
				projection.DeliveryUnitLabels().String(): "",
			},
			DeliveryUnit: &domain.DeliveryUnitFilter{
				Labels: []domain.LabelFilter{
					{Type: "NON_EXISTENT", Value: "VALUE"},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
		Expect(hasMore).To(BeFalse())
	})

	It("should filter delivery units by LPNs", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear delivery units con diferentes LPNs
		deliveryUnit1 := domain.DeliveryUnit{
			Lpn: "LPN001",
		}

		deliveryUnit2 := domain.DeliveryUnit{
			Lpn: "LPN002",
		}

		deliveryUnit3 := domain.DeliveryUnit{
			Lpn: "LPN003",
		}

		// Crear orders con las delivery units
		orders := []domain.Order{
			{
				ReferenceID: "REF-001",
				DeliveryUnits: []domain.DeliveryUnit{
					deliveryUnit1,
				},
			},
			{
				ReferenceID: "REF-002",
				DeliveryUnits: []domain.DeliveryUnit{
					deliveryUnit2,
				},
			},
			{
				ReferenceID: "REF-003",
				DeliveryUnits: []domain.DeliveryUnit{
					deliveryUnit3,
				},
			},
		}

		// Insertar delivery units
		err = NewUpsertDeliveryUnits(conn)(ctx, []domain.DeliveryUnit{deliveryUnit1, deliveryUnit2, deliveryUnit3})
		Expect(err).ToNot(HaveOccurred())

		// Insertar orders
		for _, order := range orders {
			err = NewUpsertOrder(conn)(ctx, order)
			Expect(err).ToNot(HaveOccurred())
		}

		// Crear delivery units history
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: orders,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			projection)

		// Test filtrando por un solo LPN
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():     "",
				projection.DeliveryUnitLPN().String(): "",
			},
			DeliveryUnit: &domain.DeliveryUnitFilter{
				Lpns: []string{"LPN001"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())
		Expect(results[0].OrderReferenceID).To(Equal("REF-001"))

		// Test filtrando por múltiples LPNs
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():     "",
				projection.DeliveryUnitLPN().String(): "",
			},
			DeliveryUnit: &domain.DeliveryUnitFilter{
				Lpns: []string{"LPN001", "LPN003"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(2))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-001")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-003")),
		))

		// Test filtrando por LPN no existente
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():     "",
				projection.DeliveryUnitLPN().String(): "",
			},
			DeliveryUnit: &domain.DeliveryUnitFilter{
				Lpns: []string{"NON-EXISTENT-LPN"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
		Expect(hasMore).To(BeFalse())
	})

	It("should filter delivery units by coordinates confidence level according to the test matrix", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear direcciones con diferentes niveles de confianza según la matriz de prueba
		destinations := []domain.AddressInfo{
			{
				State:        "CA",
				Province:     "CA",
				District:     "CA",
				AddressLine1: "Address A",
				Coordinates: domain.Coordinates{
					Point:  orb.Point{1, 1},
					Source: "geocoding",
					Confidence: domain.CoordinatesConfidence{
						Level:   0.3,
						Message: "Low confidence",
						Reason:  "Geocoding service",
					},
				},
			},
			{
				State:        "CA",
				Province:     "CA",
				District:     "CA",
				AddressLine1: "Address B",
				Coordinates: domain.Coordinates{
					Point:  orb.Point{2, 2},
					Source: "geocoding",
					Confidence: domain.CoordinatesConfidence{
						Level:   0.4,
						Message: "Low confidence",
						Reason:  "Geocoding service",
					},
				},
			},
			{
				State:        "CA",
				Province:     "CA",
				District:     "CA",
				AddressLine1: "Address C",
				Coordinates: domain.Coordinates{
					Point:  orb.Point{3, 3},
					Source: "geocoding",
					Confidence: domain.CoordinatesConfidence{
						Level:   0.6,
						Message: "Medium confidence",
						Reason:  "Geocoding service",
					},
				},
			},
			{
				State:        "CA",
				Province:     "CA",
				District:     "CA",
				AddressLine1: "Address D",
				Coordinates: domain.Coordinates{
					Point:  orb.Point{4, 4},
					Source: "geocoding",
					Confidence: domain.CoordinatesConfidence{
						Level:   0.7,
						Message: "Medium confidence",
						Reason:  "Geocoding service",
					},
				},
			},
			{
				State:        "CA",
				Province:     "CA",
				District:     "CA",
				AddressLine1: "Address E",
				Coordinates: domain.Coordinates{
					Point:  orb.Point{5, 5},
					Source: "geocoding",
					Confidence: domain.CoordinatesConfidence{
						Level:   0.9,
						Message: "High confidence",
						Reason:  "Geocoding service",
					},
				},
			},
			{
				State:        "CA",
				Province:     "CA",
				District:     "CA",
				AddressLine1: "Address F",
				Coordinates: domain.Coordinates{
					Point:  orb.Point{6, 6},
					Source: "geocoding",
					Confidence: domain.CoordinatesConfidence{
						Level:   1.0,
						Message: "Very high confidence",
						Reason:  "Geocoding service",
					},
				},
			},
		}

		// Insertar direcciones
		for _, dest := range destinations {
			err = NewUpsertAddressInfo(conn)(ctx, dest)
			Expect(err).ToNot(HaveOccurred())
		}

		// Crear orders con las direcciones
		orders := make([]domain.Order, len(destinations))
		for i, dest := range destinations {
			orders[i] = domain.Order{
				ReferenceID: domain.ReferenceID(fmt.Sprintf("REF-%c", 'A'+i)),
				Destination: domain.NodeInfo{
					AddressInfo: dest,
				},
				DeliveryUnits: []domain.DeliveryUnit{
					{},
				},
			}
		}

		// Insertar orders
		for _, order := range orders {
			err = NewUpsertOrder(conn)(ctx, order)
			Expect(err).ToNot(HaveOccurred())
		}

		// Crear delivery units history
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: orders,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()
		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			projection)

		// Caso 1: Sin filtros
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.DestinationAddressInfo().String():                "",
				projection.ReferenceID().String():                           "",
				projection.DestinationCoordinatesConfidenceLevel().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(6))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-A")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-B")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-C")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-D")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-E")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-F")),
		))

		// Caso 2: Min = 0.5, Max = nil
		minConfidence := 0.5
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                           "",
				projection.DestinationCoordinatesConfidenceLevel().String(): "",
			},
			Destination: &domain.LocationFilter{
				CoordinatesConfidence: &domain.CoordinatesConfidenceLevelFilter{
					Min: &minConfidence,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(4))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-C")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-D")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-E")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-F")),
		))

		// Caso 3: Min = nil, Max = 0.5
		maxConfidence := 0.5
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                           "",
				projection.DestinationCoordinatesConfidenceLevel().String(): "",
			},
			Destination: &domain.LocationFilter{
				CoordinatesConfidence: &domain.CoordinatesConfidenceLevelFilter{
					Max: &maxConfidence,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(2))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-A")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-B")),
		))

		// Caso 4: Min = 0.5, Max = 0.8
		minConfidence = 0.5
		maxConfidence = 0.8
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                           "",
				projection.DestinationCoordinatesConfidenceLevel().String(): "",
			},
			Destination: &domain.LocationFilter{
				CoordinatesConfidence: &domain.CoordinatesConfidenceLevelFilter{
					Min: &minConfidence,
					Max: &maxConfidence,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(2))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-C")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-D")),
		))

		// Caso 5: Min = 0.7, Max = 0.7 (match exacto)
		minConfidence = 0.7
		maxConfidence = 0.7
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                           "",
				projection.DestinationCoordinatesConfidenceLevel().String(): "",
			},
			Destination: &domain.LocationFilter{
				CoordinatesConfidence: &domain.CoordinatesConfidenceLevelFilter{
					Min: &minConfidence,
					Max: &maxConfidence,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results[0].OrderReferenceID).To(Equal("REF-D"))
		Expect(hasMore).To(BeFalse())

		// Caso 6: Min = 1.1, Max = nil (ningún registro cumple)
		minConfidence = 1.1
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                           "",
				projection.DestinationCoordinatesConfidenceLevel().String(): "",
			},
			Destination: &domain.LocationFilter{
				CoordinatesConfidence: &domain.CoordinatesConfidenceLevelFilter{
					Min: &minConfidence,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
		Expect(hasMore).To(BeFalse())

		// Caso 7: Min = nil, Max = 0.0 (ningún registro cumple)
		maxConfidence = 0.0
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                           "",
				projection.DestinationCoordinatesConfidenceLevel().String(): "",
			},
			Destination: &domain.LocationFilter{
				CoordinatesConfidence: &domain.CoordinatesConfidenceLevelFilter{
					Max: &maxConfidence,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
		Expect(hasMore).To(BeFalse())

		// Caso 8: Min = 0.0, Max = 1.0 (todos los registros)
		minConfidence = 0.0
		maxConfidence = 1.0
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                           "",
				projection.DestinationCoordinatesConfidenceLevel().String(): "",
			},
			Destination: &domain.LocationFilter{
				CoordinatesConfidence: &domain.CoordinatesConfidenceLevelFilter{
					Min: &minConfidence,
					Max: &maxConfidence,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(6))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-A")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-B")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-C")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-D")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-E")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-F")),
		))
	})

	It("should filter delivery units by origin node references", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear nodo de origen
		originNode := domain.NodeInfo{
			ReferenceID: "NODE-001",
		}
		err = NewUpsertNodeInfo(conn)(ctx, originNode)
		Expect(err).ToNot(HaveOccurred())

		// Crear orden con el nodo de origen
		order := domain.Order{
			ReferenceID: "ORDER-001",
			Origin:      originNode,
			DeliveryUnits: []domain.DeliveryUnit{
				{},
			},
		}
		err = NewUpsertOrder(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Crear otra orden con un nodo de origen diferente
		otherOriginNode := domain.NodeInfo{
			ReferenceID: "NODE-002",
		}
		err = NewUpsertNodeInfo(conn)(ctx, otherOriginNode)
		Expect(err).ToNot(HaveOccurred())

		otherOrder := domain.Order{
			ReferenceID: "ORDER-002",
			Origin:      otherOriginNode,
			DeliveryUnits: []domain.DeliveryUnit{
				{},
			},
		}
		err = NewUpsertOrder(conn)(ctx, otherOrder)
		Expect(err).ToNot(HaveOccurred())

		// Crear historial de unidades de entrega
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: []domain.Order{order, otherOrder},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		// Buscar unidades de entrega filtrando por el primer nodo de origen
		projection := deliveryunits.NewProjection()
		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			projection,
		)

		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
			Origin: &domain.LocationFilter{
				NodeReferences: []string{"NODE-001"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1), "Debería retornar solo una unidad de entrega")
		Expect(hasMore).To(BeFalse())
		Expect(results[0].OrderReferenceID).To(Equal("ORDER-001"), "Debería retornar la orden con el nodo de origen especificado")

		// Buscar unidades de entrega filtrando por el segundo nodo de origen
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
			Origin: &domain.LocationFilter{
				NodeReferences: []string{"NODE-002"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1), "Debería retornar solo una unidad de entrega")
		Expect(hasMore).To(BeFalse())
		Expect(results[0].OrderReferenceID).To(Equal("ORDER-002"), "Debería retornar la orden con el nodo de origen especificado")

		// Buscar unidades de entrega filtrando por ambos nodos de origen
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			Origin: &domain.LocationFilter{
				NodeReferences: []string{"NODE-001", "NODE-002"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(2), "Debería retornar ambas unidades de entrega")
		Expect(hasMore).To(BeFalse())
	})

	It("should filter by onlyLatestStatus flag", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear order
		order := domain.Order{
			ReferenceID: "123",
			DeliveryUnits: []domain.DeliveryUnit{
				{},
			},
		}
		err = NewUpsertOrder(conn)(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Crear múltiples estados para la misma unidad de entrega
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

		// Crear otro estado para la misma unidad
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

		// Test con onlyLatestStatus = true
		resultsWithLatest, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			OnlyLatestStatus: true,
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(resultsWithLatest).To(HaveLen(1), "Debería retornar solo el último estado")
		Expect(hasMore).To(BeFalse())

		// Test con onlyLatestStatus = false
		resultsWithoutLatest, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			OnlyLatestStatus: false,
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(resultsWithoutLatest).To(HaveLen(2), "Debería retornar todos los estados")
		Expect(hasMore).To(BeFalse())

		// Verificar que el ID del resultado con onlyLatestStatus es mayor que el otro
		Expect(resultsWithLatest[0].ID).To(BeNumerically(">", resultsWithoutLatest[0].ID),
			"El ID del último estado debería ser mayor")
	})

	It("should filter delivery units by promised date range", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear dos órdenes con diferentes fechas prometidas
		fixedDate1 := time.Date(2025, 5, 26, 0, 0, 0, 0, time.UTC)
		fixedDate2 := time.Date(2025, 5, 28, 0, 0, 0, 0, time.UTC)

		// Primera orden
		order1 := domain.Order{
			ReferenceID: "123",
			Destination: domain.NodeInfo{
				AddressInfo: domain.AddressInfo{
					State:    "CA",
					Province: "CA",
					District: "CA",
				},
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{},
			},
			PromisedDate: domain.PromisedDate{
				DateRange: domain.DateRange{
					StartDate: fixedDate1,
					EndDate:   fixedDate1.Add(24 * time.Hour),
				},
				TimeRange: domain.TimeRange{
					StartTime: "10:00",
					EndTime:   "17:00",
				},
				ServiceCategory: "STANDARD",
			},
		}
		err = NewUpsertOrder(conn)(ctx, order1)
		Expect(err).ToNot(HaveOccurred())

		// Segunda orden
		order2 := domain.Order{
			ReferenceID: "456",
			Destination: domain.NodeInfo{
				AddressInfo: domain.AddressInfo{
					State:    "CA",
					Province: "CA",
					District: "CA",
				},
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{},
			},
			PromisedDate: domain.PromisedDate{
				DateRange: domain.DateRange{
					StartDate: fixedDate2,
					EndDate:   fixedDate2.Add(24 * time.Hour),
				},
				TimeRange: domain.TimeRange{
					StartTime: "10:00",
					EndTime:   "17:00",
				},
				ServiceCategory: "STANDARD",
			},
		}
		err = NewUpsertOrder(conn)(ctx, order2)
		Expect(err).ToNot(HaveOccurred())

		// Crear historial de unidades de entrega
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: []domain.Order{order1, order2},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		// Probar filtro con rango que incluye solo la primera orden
		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			deliveryunits.NewProjection())

		startDate := "2025-05-25T00:00:00Z"
		endDate := "2025-05-27T00:00:00Z"

		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                    "",
				projection.PromisedDateDateRangeStartDate().String(): "",
				projection.PromisedDateDateRangeEndDate().String():   "",
			},
			PromisedDate: &domain.PromisedDateFilter{
				DateRange: &domain.DateRangeFilter{
					StartDate: &startDate,
					EndDate:   &endDate,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1), "Expected 1 result, got %d", len(results))
		Expect(hasMore).To(BeFalse())
		Expect(results[0].OrderReferenceID).To(Equal("123"), "Unexpected reference ID")
		Expect(results[0].OrderPromisedDateStartDate).To(Equal("2025-05-26T00:00:00Z"), "Unexpected promised date start")
		Expect(results[0].OrderPromisedDateEndDate).To(Equal("2025-05-27T00:00:00Z"), "Unexpected promised date end")

		// Probar filtro con rango que incluye ambas órdenes
		startDate = "2025-05-25T00:00:00Z"
		endDate = "2025-05-29T00:00:00Z"

		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():                    "",
				projection.PromisedDateDateRangeStartDate().String(): "",
				projection.PromisedDateDateRangeEndDate().String():   "",
			},
			PromisedDate: &domain.PromisedDateFilter{
				DateRange: &domain.DateRangeFilter{
					StartDate: &startDate,
					EndDate:   &endDate,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(2), "Expected 2 results, got %d", len(results))
		Expect(hasMore).To(BeFalse())

		// Verificar que las órdenes están ordenadas por reference_id
		Expect(results[0].OrderReferenceID).To(Equal("123"), "Unexpected first reference ID")
		Expect(results[1].OrderReferenceID).To(Equal("456"), "Unexpected second reference ID")
	})

	It("should filter delivery units by size categories", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear delivery units con diferentes size categories
		deliveryUnit1 := domain.DeliveryUnit{
			Lpn:          "1234567890",
			SizeCategory: domain.SizeCategory{Code: "SMALL"},
		}

		deliveryUnit2 := domain.DeliveryUnit{
			Lpn:          "1234567891",
			SizeCategory: domain.SizeCategory{Code: "MEDIUM"},
		}

		deliveryUnit3 := domain.DeliveryUnit{
			Lpn:          "1234567892",
			SizeCategory: domain.SizeCategory{Code: "LARGE"},
		}

		// Crear orders con las delivery units
		orders := []domain.Order{
			{
				ReferenceID: "REF-001",
				DeliveryUnits: []domain.DeliveryUnit{
					deliveryUnit1,
				},
			},
			{
				ReferenceID: "REF-002",
				DeliveryUnits: []domain.DeliveryUnit{
					deliveryUnit2,
				},
			},
			{
				ReferenceID: "REF-003",
				DeliveryUnits: []domain.DeliveryUnit{
					deliveryUnit3,
				},
			},
		}

		// Insertar delivery units
		err = NewUpsertDeliveryUnits(conn)(ctx, []domain.DeliveryUnit{deliveryUnit1, deliveryUnit2, deliveryUnit3})
		Expect(err).ToNot(HaveOccurred())

		// Insertar orders
		for _, order := range orders {
			err = NewUpsertOrder(conn)(ctx, order)
			Expect(err).ToNot(HaveOccurred())
		}

		// Crear delivery units history
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: orders,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			projection)

		// Test filtrando por una sola size category
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
			DeliveryUnit: &domain.DeliveryUnitFilter{
				SizeCategories: []string{"SMALL"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())
		Expect(results[0].OrderReferenceID).To(Equal("REF-001"))

		// Test filtrando por múltiples size categories
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String(): "",
			},
			DeliveryUnit: &domain.DeliveryUnitFilter{
				SizeCategories: []string{"SMALL", "LARGE"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(2))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-001")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-003")),
		))
	})

	It("should filter delivery units by collect availability dates", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear fechas fijas para las pruebas
		date1 := time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)
		date2 := time.Date(2024, 3, 21, 0, 0, 0, 0, time.UTC)
		date3 := time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC)

		// Crear orders con diferentes fechas de disponibilidad de recolección
		orders := []domain.Order{
			{
				ReferenceID: "REF-001",
				CollectAvailabilityDate: domain.CollectAvailabilityDate{
					Date: date1,
					TimeRange: domain.TimeRange{
						StartTime: "09:00",
						EndTime:   "18:00",
					},
				},
				DeliveryUnits: []domain.DeliveryUnit{
					{},
				},
			},
			{
				ReferenceID: "REF-002",
				CollectAvailabilityDate: domain.CollectAvailabilityDate{
					Date: date2,
					TimeRange: domain.TimeRange{
						StartTime: "10:00",
						EndTime:   "19:00",
					},
				},
				DeliveryUnits: []domain.DeliveryUnit{
					{},
				},
			},
			{
				ReferenceID: "REF-003",
				CollectAvailabilityDate: domain.CollectAvailabilityDate{
					Date: date3,
					TimeRange: domain.TimeRange{
						StartTime: "11:00",
						EndTime:   "20:00",
					},
				},
				DeliveryUnits: []domain.DeliveryUnit{
					{},
				},
			},
		}

		// Insertar orders
		for _, order := range orders {
			err = NewUpsertOrder(conn)(ctx, order)
			Expect(err).ToNot(HaveOccurred())
		}

		// Crear delivery units history
		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: orders,
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		projection := deliveryunits.NewProjection()

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(
			conn,
			projection)

		// Test filtrando por una sola fecha
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():             "",
				projection.CollectAvailabilityDate().String(): "",
			},
			CollectAvailability: &domain.CollectAvailabilityFilter{
				Dates: []string{"2024-03-20"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())
		Expect(results[0].OrderReferenceID).To(Equal("REF-001"))

		// Test filtrando por múltiples fechas
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():             "",
				projection.CollectAvailabilityDate().String(): "",
			},
			CollectAvailability: &domain.CollectAvailabilityFilter{
				Dates: []string{"2024-03-20", "2024-03-22"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(2))
		Expect(hasMore).To(BeFalse())
		Expect(results).To(ContainElements(
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-001")),
			WithTransform(func(r projectionresult.DeliveryUnitsProjectionResult) string {
				return r.OrderReferenceID
			}, Equal("REF-003")),
		))

		// Test filtrando por fecha no existente
		results, hasMore, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():             "",
				projection.CollectAvailabilityDate().String(): "",
			},
			CollectAvailability: &domain.CollectAvailabilityFilter{
				Dates: []string{"2024-03-25"},
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
		Expect(hasMore).To(BeFalse())
	})

	It("should return delivery units with manual change information", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear order con manual change
		order := domain.Order{
			ReferenceID: "123",
			DeliveryUnits: []domain.DeliveryUnit{
				{
					ConfirmDelivery: domain.ConfirmDelivery{
						ManualChange: domain.ManualChange{
							PerformedBy: "juan@example.com",
							Reason:      "Corrección tras reclamo de transporte",
						},
					},
				},
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
		results, hasMore, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{
			RequestedFields: map[string]any{
				projection.ReferenceID().String():             "",
				projection.ManualChange().String():            "",
				projection.ManualChangePerformedBy().String(): "",
				projection.ManualChangeReason().String():      "",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(hasMore).To(BeFalse())

		// Validar que la información de manual change se recuperó correctamente
		Expect(results[0].OrderReferenceID).To(Equal("123"))
		Expect(results[0].ManualChangePerformedBy).To(Equal("juan@example.com"))
		Expect(results[0].ManualChangeReason).To(Equal("Corrección tras reclamo de transporte"))
	})

})

// Helper function to create pointer to int
func ptr(i int) *int {
	return &i
}
