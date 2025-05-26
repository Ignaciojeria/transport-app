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
		Expect(results[0].DestinationDistrict).To(Equal("CA"))
		Expect(results[0].DestinationLatitude).To(Equal(1.0))
		Expect(results[0].DestinationLongitude).To(Equal(1.0))
		Expect(results[0].DestinationProvince).To(Equal("CA"))
		Expect(results[0].DestinationState).To(Equal("CA"))
		Expect(results[0].DestinationTimeZone).To(Equal("America/Santiago"))
		Expect(results[0].DestinationZipCode).To(Equal("12345"))
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
