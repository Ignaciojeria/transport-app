package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertOrderHeaders", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertOrderHeaders
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertOrderHeaders(conn, nil)
	})

	It("should insert headers if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		headers := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		err = upsert(ctx, headers)
		Expect(err).ToNot(HaveOccurred())

		var dbHeaders table.OrderHeaders
		err = conn.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", headers.DocID(ctx)).
			First(&dbHeaders).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbHeaders.Commerce).To(Equal("CORP"))
		Expect(dbHeaders.Consumer).To(Equal("CORP"))
		Expect(dbHeaders.Channel).To(Equal("WEB"))
		Expect(dbHeaders.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update headers if they are different", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Headers{
			Commerce: "CORP",
			Consumer: "RETAIL",
			Channel:  "WEB",
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbHeaders table.OrderHeaders
		err = conn.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbHeaders).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbHeaders.Commerce).To(Equal("CORP"))
		Expect(dbHeaders.Consumer).To(Equal("RETAIL"))
		Expect(dbHeaders.Channel).To(Equal("WEB"))
		Expect(dbHeaders.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should not update if headers are the same", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		headers := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		err = upsert(ctx, headers)
		Expect(err).ToNot(HaveOccurred())

		// Get initial timestamp
		var initialHeaders table.OrderHeaders
		err = conn.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", headers.DocID(ctx)).
			First(&initialHeaders).Error
		Expect(err).ToNot(HaveOccurred())
		initialUpdatedAt := initialHeaders.UpdatedAt

		// Try to update with same values
		err = upsert(ctx, headers)
		Expect(err).ToNot(HaveOccurred())

		var dbHeaders table.OrderHeaders
		err = conn.DB.WithContext(ctx).
			Table("order_headers").
			Where("document_id = ?", headers.DocID(ctx)).
			First(&dbHeaders).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbHeaders.UpdatedAt).To(Equal(initialUpdatedAt))
		Expect(dbHeaders.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should allow same headers for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		headers1 := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		headers2 := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		err = upsert(ctx1, headers1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, headers2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := headers1.DocID(ctx1)
		docID2 := headers2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each header belongs to its respective tenant using DocID
		var dbHeaders1, dbHeaders2 table.OrderHeaders
		err = conn.DB.WithContext(ctx1).
			Table("order_headers").
			Where("document_id = ?", docID1).
			First(&dbHeaders1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbHeaders1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("order_headers").
			Where("document_id = ?", docID2).
			First(&dbHeaders2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbHeaders2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no order_headers table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		headers := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		upsert := NewUpsertOrderHeaders(noTablesContainerConnection, nil)
		err = upsert(ctx, headers)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("order_headers"))
	})
})
