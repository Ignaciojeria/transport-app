package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertOrderPackages", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should insert one package with LPN", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-LPN-001",
			Headers:     domain.Headers{Commerce: "c1", Consumer: "c2"},
			Packages: []domain.Package{
				{
					Lpn: "PKG-001",
					Items: []domain.Item{
						{Sku: "ITEM-1"},
					},
				},
			},
		}

		uop := NewUpsertOrderDeliveryUnits(conn)
		err = uop(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = conn.DB.WithContext(ctx).
			Table("order_delivery_units").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))
	})

	It("should replace old packages with new ones", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-REPLACE-001",
			Headers:     domain.Headers{Commerce: "c1", Consumer: "c2"},
			Packages: []domain.Package{
				{Lpn: "PKG-OLD", Items: []domain.Item{{Sku: "OLD"}}},
			},
		}

		uop := NewUpsertOrderDeliveryUnits(conn)
		err = uop(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		updated := order
		updated.Packages = []domain.Package{
			{Lpn: "PKG-NEW", Items: []domain.Item{{Sku: "NEW"}}},
		}
		err = uop(ctx, updated)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = conn.DB.WithContext(ctx).
			Table("order_delivery_units").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))
	})

	It("should not fail if no packages are provided", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-EMPTY",
			Headers:     domain.Headers{Commerce: "x", Consumer: "y"},
			Packages:    []domain.Package{},
		}

		uop := NewUpsertOrderDeliveryUnits(conn)
		err = uop(ctx, order)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail if the order_packages table is missing", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-FAIL",
			Packages: []domain.Package{
				{Lpn: "PKG-ERR", Items: []domain.Item{{Sku: "E"}}},
			},
		}

		uop := NewUpsertOrderDeliveryUnits(noTablesContainerConnection)
		err = uop(ctx, order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("order_delivery_units"))
	})

	It("should insert placeholder package with a valid docID when no packages are provided", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-PLACEHOLDER-001",
			Headers:     domain.Headers{Commerce: "a", Consumer: "b"},
			Packages:    []domain.Package{}, // no paquetes
		}

		uop := NewUpsertOrderDeliveryUnits(conn)
		err = uop(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var results []table.OrderDeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("order_delivery_units").
			Where("order_doc = ? ", order.DocID(ctx)).
			Find(&results).Error
		Expect(err).ToNot(HaveOccurred())

		Expect(results).To(HaveLen(1))
		Expect(results[0].DeliveryUnitDoc).ToNot(BeEmpty()) // docID generado desde pkg vacío + refID
	})
})
