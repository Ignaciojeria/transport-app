package tidbrepository

import (
	"context"
	"time"
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
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

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
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(results[0].ID).To(Equal(int64(1)))
		Expect(results[0].OrderReferenceID).To(Equal("123"))
		Expect(results[0].OrderCollectAvailabilityDate).To(Equal("2025-05-26T00:00:00Z"))
		Expect(results[0].OrderCollectAvailabilityDateStartTime).To(Equal("09:00"))
		Expect(results[0].OrderCollectAvailabilityDateEndTime).To(Equal("18:00"))
		Expect(results[0].OrderPromisedDateStartDate).To(Equal("2025-05-26T00:00:00Z"))
		Expect(results[0].OrderPromisedDateEndDate).To(Equal("2025-05-27T00:00:00Z"))
		Expect(results[0].OrderPromisedDateStartTime).To(Equal("10:00"))
		Expect(results[0].OrderPromisedDateEndTime).To(Equal("17:00"))
		Expect(results[0].OrderPromisedDateServiceCategory).To(Equal("STANDARD"))

		// Validaciones de Destination Address
		Expect(results[0].DestinationAddressLine1).To(Equal("123 Main St"))
		Expect(results[0].DestinationAddressLine2).To(Equal("Apt 1"))
		Expect(results[0].DestinationDistrict).To(Equal("CA"))
		Expect(results[0].DestinationLatitude).To(Equal(1.0))
		Expect(results[0].DestinationLongitude).To(Equal(1.0))
		Expect(results[0].DestinationProvince).To(Equal("CA"))
		Expect(results[0].DestinationState).To(Equal("CA"))
		Expect(results[0].DestinationTimeZone).To(Equal("America/Santiago"))
		Expect(results[0].DestinationZipCode).To(Equal("12345"))
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
		Expect(results[0].DestinationAdditionalContactMethods).To(ContainSubstring(`"type":"whatsapp"`))
		Expect(results[0].DestinationAdditionalContactMethods).To(ContainSubstring(`"value":"+56987654321"`))
		Expect(results[0].DestinationAdditionalContactMethods).To(ContainSubstring(`"type":"telegram"`))
		Expect(results[0].DestinationAdditionalContactMethods).To(ContainSubstring(`"value":"@johndoe"`))
		Expect(results[0].DestinationContactDocuments).To(ContainSubstring(`"type":"dni"`))
		Expect(results[0].DestinationContactDocuments).To(ContainSubstring(`"value":"12345678-9"`))
		Expect(results[0].DestinationContactDocuments).To(ContainSubstring(`"type":"passport"`))
		Expect(results[0].DestinationContactDocuments).To(ContainSubstring(`"value":"AB123456"`))
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

		err = NewUpsertDeliveryUnits(conn)(ctx, []domain.DeliveryUnit{deliveryUnit}, "123")
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
		Expect(results[0].LPN).To(Equal("LPN123"))
		Expect(results[0].JSONDimensions).To(ContainSubstring(`"length":10`))
		Expect(results[0].JSONDimensions).To(ContainSubstring(`"width":20`))
		Expect(results[0].JSONDimensions).To(ContainSubstring(`"height":30`))
		Expect(results[0].JSONDimensions).To(ContainSubstring(`"unit":"cm"`))
		Expect(results[0].JSONWeight).To(ContainSubstring(`"weight_value":5.5`))
		Expect(results[0].JSONWeight).To(ContainSubstring(`"weight_unit":"kg"`))
		Expect(results[0].JSONInsurance).To(ContainSubstring(`"currency":"USD"`))
		Expect(results[0].JSONInsurance).To(ContainSubstring(`"unit_value":100`))
		Expect(results[0].JSONItems).To(ContainSubstring(`"sku":"SKU123"`))
		Expect(results[0].JSONItems).To(ContainSubstring(`"description":"Test Item"`))
		Expect(results[0].JSONItems).To(ContainSubstring(`"quantity_number":2`))
		Expect(results[0].JSONItems).To(ContainSubstring(`"quantity_unit":"pcs"`))

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
})
