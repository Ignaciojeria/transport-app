package tidbrepository

import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertOrderReferences", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should insert order references correctly", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-REF-001",
			Headers: domain.Headers{
				Commerce: "commerce-1",
				Consumer: "consumer-1",
			},
			References: []domain.Reference{
				{Type: "TRACKING", Value: "123456"},
				{Type: "EXTERNAL", Value: "ABCDEF"},
			},
		}

		uor := NewUpsertOrderReferences(conn, nil)
		err = uor(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = conn.DB.WithContext(ctx).
			Table("order_references").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	It("should replace old references with new ones", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-REF-002",
			Headers:     domain.Headers{Commerce: "c", Consumer: "c"},
			References: []domain.Reference{
				{Type: "OLD", Value: "1"},
			},
		}

		uor := NewUpsertOrderReferences(conn, nil)
		err = uor(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		updated := order
		updated.References = []domain.Reference{
			{Type: "NEW", Value: "2"},
		}
		err = uor(ctx, updated)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = conn.DB.WithContext(ctx).
			Table("order_references").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))
	})

	It("should not fail if no references are provided", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-REF-003",
			Headers:     domain.Headers{Commerce: "x", Consumer: "y"},
			References:  []domain.Reference{},
		}

		uor := NewUpsertOrderReferences(conn, nil)
		err = uor(ctx, order)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail if the order_references table is missing", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		uor := NewUpsertOrderReferences(noTablesContainerConnection, nil)
		order := domain.Order{
			ReferenceID: "ORD-REF-ERR",
			References: []domain.Reference{
				{Type: "ERR", Value: "ERR"},
			},
		}
		err = uor(ctx, order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("order_references"))
	})

	It("should insert placeholder when no references are provided", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		order := domain.Order{
			ReferenceID: "ORD-REF-PLACEHOLDER-001",
			Headers:     domain.Headers{Commerce: "z", Consumer: "w"},
			References:  []domain.Reference{},
		}

		uor := NewUpsertOrderReferences(conn, nil)
		err = uor(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = conn.DB.WithContext(ctx).
			Table("order_references").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))
	})
})
