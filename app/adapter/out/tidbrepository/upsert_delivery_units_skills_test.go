package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertDeliveryUnitsSkills", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertDeliveryUnitsSkills
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertDeliveryUnitsSkills(conn)
	})

	It("should insert delivery unit skills correctly", func() {
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
					Skills: []domain.Skill{
						"DRIVING",
						"LIFTING",
					},
				},
			},
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_skills").
			Where("delivery_unit_doc = ?", order.DeliveryUnits[0].DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify the skills were saved correctly
		var skills []table.DeliveryUnitsSkills
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_skills").
			Where("delivery_unit_doc = ?", order.DeliveryUnits[0].DocID(ctx)).
			Find(&skills).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(skills).To(HaveLen(2))
	})

	It("should update existing skills", func() {
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
					Skills: []domain.Skill{
						"DRIVING",
						"LIFTING",
					},
				},
			},
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Update with new skills
		order.DeliveryUnits[0].Skills = []domain.Skill{
			"DRIVING",
			"PACKAGING",
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Verify the skills were updated
		var skills []table.DeliveryUnitsSkills
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_skills").
			Where("delivery_unit_doc = ?", order.DeliveryUnits[0].DocID(ctx)).
			Find(&skills).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(skills).To(HaveLen(2))
	})

	It("should handle empty skills", func() {
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
					Skills: []domain.Skill{},
				},
			},
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_skills").
			Where("delivery_unit_doc = ?", order.DeliveryUnits[0].DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))
	})

	It("should delete old skills when updating", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// First insert with multiple skills
		order := domain.Order{
			ReferenceID: "ORD-REF-004",
			Headers: domain.Headers{
				Commerce: "commerce-1",
				Consumer: "consumer-1",
			},
			DeliveryUnits: []domain.DeliveryUnit{
				{
					Lpn: "LPN-004",
					Skills: []domain.Skill{
						"DRIVING",
						"LIFTING",
						"PACKAGING",
					},
				},
			},
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Update with fewer skills
		order.DeliveryUnits[0].Skills = []domain.Skill{
			"DRIVING",
		}

		err = upsert(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		// Verify only one skill remains
		var count int64
		err = conn.DB.WithContext(ctx).
			Table("delivery_units_skills").
			Where("delivery_unit_doc = ?", order.DeliveryUnits[0].DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))
	})
})
