package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertDeliveryUnitsLabels", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertDeliveryUnitsLabels
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertDeliveryUnitsLabels(conn)
	})

	It("should insert delivery unit labels correctly", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-REF-001",
			Headers: domain.Headers{
				Commerce: "commerce-1",
				Consumer: "consumer-1",
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{
					Lpn: "LPN-001",
					Labels: []domain.Reference{
						{Type: "TRACKING", Value: "123456"},
						{Type: "EXTERNAL", Value: "ABCDEF"},
					},
				},
			},
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_labels").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify the labels were saved correctly
		var labels []table.DeliveryUnitsLabels
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_labels").
			Where("order_doc = ?", order.DocID(ctx)).
			Find(&labels).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(labels).To(HaveLen(2))

		// Verify each label
		labelMap := make(map[string]table.DeliveryUnitsLabels)
		for _, label := range labels {
			labelMap[label.Type] = label
		}

		Expect(labelMap["TRACKING"].Value).To(Equal("123456"))
		Expect(labelMap["EXTERNAL"].Value).To(Equal("ABCDEF"))
	})

	It("should update existing labels", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// First insert
		order := domain.Order{
			ReferenceID: "ORD-REF-002",
			Headers: domain.Headers{
				Commerce: "commerce-1",
				Consumer: "consumer-1",
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{
					Lpn: "LPN-002",
					Labels: []domain.Reference{
						{Type: "TRACKING", Value: "123456"},
						{Type: "EXTERNAL", Value: "ABCDEF"},
					},
				},
			},
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Update with new labels
		order.DeliveryUnits[0].Labels = []domain.Reference{
			{Type: "TRACKING", Value: "654321"},
			{Type: "NEW_TYPE", Value: "NEW_VALUE"},
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Verify the labels were updated
		var labels []table.DeliveryUnitsLabels
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_labels").
			Where("order_doc = ?", order.DocID(ctx)).
			Find(&labels).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(labels).To(HaveLen(2))

		// Verify each label
		labelMap := make(map[string]table.DeliveryUnitsLabels)
		for _, label := range labels {
			labelMap[label.Type] = label
		}

		Expect(labelMap["TRACKING"].Value).To(Equal("654321"))
		Expect(labelMap["NEW_TYPE"].Value).To(Equal("NEW_VALUE"))
	})

	It("should handle empty labels", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-REF-003",
			Headers: domain.Headers{
				Commerce: "commerce-1",
				Consumer: "consumer-1",
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{
					Lpn:    "LPN-003",
					Labels: []domain.Reference{},
				},
			},
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_labels").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))

		// Verify the empty label was created
		var label table.DeliveryUnitsLabels
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_labels").
			Where("order_doc = ?", order.DocID(ctx)).
			First(&label).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(label.Type).To(Equal(""))
		Expect(label.Value).To(Equal(""))
	})

	It("should delete old labels when updating", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// First insert with multiple labels
		order := domain.Order{
			ReferenceID: "ORD-REF-004",
			Headers: domain.Headers{
				Commerce: "commerce-1",
				Consumer: "consumer-1",
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{
					Lpn: "LPN-004",
					Labels: []domain.Reference{
						{Type: "TRACKING", Value: "123456"},
						{Type: "EXTERNAL", Value: "ABCDEF"},
						{Type: "EXTRA", Value: "EXTRA_VALUE"},
					},
				},
			},
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Update with fewer labels
		order.DeliveryUnits[0].Labels = []domain.Reference{
			{Type: "TRACKING", Value: "654321"},
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Verify only one label remains
		var count int64
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_labels").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))

		// Verify the remaining label
		var label table.DeliveryUnitsLabels
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_labels").
			Where("order_doc = ?", order.DocID(ctx)).
			First(&label).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(label.Type).To(Equal("TRACKING"))
		Expect(label.Value).To(Equal("654321"))
	})
})
