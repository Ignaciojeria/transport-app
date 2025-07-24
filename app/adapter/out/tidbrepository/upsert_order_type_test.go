package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertOrderType", func() {
	var (
		upsert UpsertOrderType
	)

	BeforeEach(func() {
		upsert = NewUpsertOrderType(connection, nil)
	})

	It("should insert order type if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())

		orderType := domain.OrderType{
			Type:        "TEST",
			Description: "Test Order Type",
		}

		err = upsert(ctx, orderType)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := orderType.DocID(ctx)

		var dbOrderType table.OrderType
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", docID).
			First(&dbOrderType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbOrderType.Type).To(Equal("TEST"))
		Expect(dbOrderType.Description).To(Equal("Test Order Type"))
		Expect(dbOrderType.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update existing record when description changes", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())

		original := domain.OrderType{
			Type:        "TEST",
			Description: "Test Order Type",
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Get the original record
		var originalRecord table.OrderType
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", original.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(originalRecord.TenantID.String()).To(Equal(tenant.ID.String()))

		modified := domain.OrderType{
			Type:        "TEST",
			Description: "Updated Test Order Type",
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		// Verify a new record was created
		var dbOrderType table.OrderType
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbOrderType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbOrderType.Type).To(Equal("TEST"))
		Expect(dbOrderType.Description).To(Equal("Updated Test Order Type"))
		Expect(dbOrderType.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should not update if order type is the same", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())

		orderType := domain.OrderType{
			Type:        "TEST",
			Description: "Test Order Type",
		}

		err = upsert(ctx, orderType)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := orderType.DocID(ctx)

		// Get initial timestamp
		var initialOrderType table.OrderType
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", docID).
			First(&initialOrderType).Error
		Expect(err).ToNot(HaveOccurred())
		initialUpdatedAt := initialOrderType.UpdatedAt

		// Try to update with same values
		err = upsert(ctx, orderType)
		Expect(err).ToNot(HaveOccurred())

		var dbOrderType table.OrderType
		err = connection.DB.WithContext(ctx).
			Table("order_types").
			Where("document_id = ?", docID).
			First(&dbOrderType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbOrderType.UpdatedAt).To(Equal(initialUpdatedAt))
		Expect(dbOrderType.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should allow same order type for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())

		orderType1 := domain.OrderType{
			Type:        "multi-org",
			Description: "Testing multiple organizations",
		}
		orderType2 := domain.OrderType{
			Type:        "multi-org",
			Description: "Testing multiple organizations",
		}

		err = upsert(ctx1, orderType1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, orderType2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := orderType1.DocID(ctx1)
		docID2 := orderType2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each order type belongs to its respective tenant using DocID
		var dbOrderType1, dbOrderType2 table.OrderType
		err = connection.DB.WithContext(ctx1).
			Table("order_types").
			Where("document_id = ?", docID1).
			First(&dbOrderType1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbOrderType1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = connection.DB.WithContext(ctx2).
			Table("order_types").
			Where("document_id = ?", docID2).
			First(&dbOrderType2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbOrderType2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no order_types table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), connection)
		Expect(err).ToNot(HaveOccurred())

		orderType := domain.OrderType{
			Type:        "TEST",
			Description: "Test Order Type",
		}

		upsert := NewUpsertOrderType(noTablesContainerConnection, nil)
		err = upsert(ctx, orderType)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("order_types"))
	})
})
