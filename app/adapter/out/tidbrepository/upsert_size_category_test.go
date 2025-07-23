package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertSizeCategory", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should insert size category if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		sizeCategory := domain.SizeCategory{
			Code: "SMALL",
		}

		upsert := NewUpsertSizeCategory(conn, nil)
		err = upsert(ctx, sizeCategory)
		Expect(err).ToNot(HaveOccurred())

		// Verify using document_id
		var dbSizeCategory table.SizeCategory
		err = conn.DB.WithContext(ctx).
			Table("size_categories").
			Where("document_id = ? AND tenant_id = ?", sizeCategory.DocumentID(ctx), tenant.ID).
			First(&dbSizeCategory).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbSizeCategory.Code).To(Equal("SMALL"))
		Expect(dbSizeCategory.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should not insert if size category already exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		sizeCategory := domain.SizeCategory{
			Code: "MEDIUM",
		}

		upsert := NewUpsertSizeCategory(conn, nil)
		err = upsert(ctx, sizeCategory)
		Expect(err).ToNot(HaveOccurred())

		// Get the original record
		var originalRecord table.SizeCategory
		err = conn.DB.WithContext(ctx).
			Table("size_categories").
			Where("document_id = ? AND tenant_id = ?", sizeCategory.DocumentID(ctx), tenant.ID).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Try to insert the same size category again
		err = upsert(ctx, sizeCategory)
		Expect(err).ToNot(HaveOccurred())

		// Verify the record hasn't changed
		var dbSizeCategory table.SizeCategory
		err = conn.DB.WithContext(ctx).
			Table("size_categories").
			Where("document_id = ? AND tenant_id = ?", sizeCategory.DocumentID(ctx), tenant.ID).
			First(&dbSizeCategory).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbSizeCategory.Code).To(Equal("MEDIUM"))
		Expect(dbSizeCategory.TenantID.String()).To(Equal(tenant.ID.String()))
		Expect(dbSizeCategory.UpdatedAt).To(Equal(originalRecord.UpdatedAt))
	})

	It("should allow same size category for different tenants", func() {
		// Create two different tenants
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		sizeCategory := domain.SizeCategory{
			Code: "LARGE",
		}

		upsert := NewUpsertSizeCategory(conn, nil)

		// Insert for first tenant
		err = upsert(ctx1, sizeCategory)
		Expect(err).ToNot(HaveOccurred())

		// Insert for second tenant
		err = upsert(ctx2, sizeCategory)
		Expect(err).ToNot(HaveOccurred())

		// Verify both records exist
		var dbSizeCategory1 table.SizeCategory
		err = conn.DB.WithContext(ctx1).
			Table("size_categories").
			Where("document_id = ? AND tenant_id = ?", sizeCategory.DocumentID(ctx1), tenant1.ID).
			First(&dbSizeCategory1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbSizeCategory1.Code).To(Equal("LARGE"))
		Expect(dbSizeCategory1.TenantID.String()).To(Equal(tenant1.ID.String()))

		var dbSizeCategory2 table.SizeCategory
		err = conn.DB.WithContext(ctx2).
			Table("size_categories").
			Where("document_id = ? AND tenant_id = ?", sizeCategory.DocumentID(ctx2), tenant2.ID).
			First(&dbSizeCategory2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbSizeCategory2.Code).To(Equal("LARGE"))
		Expect(dbSizeCategory2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no size_categories table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		sizeCategory := domain.SizeCategory{
			Code: "ERROR",
		}

		upsert := NewUpsertSizeCategory(noTablesContainerConnection, nil)
		err = upsert(ctx, sizeCategory)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("size_categories"))
	})
})
