package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(conn)
		results, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(BeEmpty())
	})

	It("should return delivery units when they exist", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		err = NewUpsertDeliveryUnitsHistory(conn)(ctx, domain.Plan{
			Routes: []domain.Route{
				{
					Orders: []domain.Order{
						{
							DeliveryUnits: []domain.DeliveryUnit{
								{},
							},
						},
					},
				},
			},
		})
		Expect(err).ToNot(HaveOccurred())

		err = NewUpsertOrder(conn)(ctx, domain.Order{})
		Expect(err).ToNot(HaveOccurred())

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(conn)
		results, err := findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{})
		Expect(err).ToNot(HaveOccurred())
		Expect(results).To(HaveLen(1))
		Expect(results[0].ID).To(Equal(int64(1)))
	})

	It("should fail if database has no delivery_units_histories table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		findDeliveryUnits := NewFindDeliveryUnitsProjectionResult(noTablesContainerConnection)
		_, err = findDeliveryUnits(ctx, domain.DeliveryUnitsFilter{})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("delivery_units_histories"))
	})
})
