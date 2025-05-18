package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertNodeInfoHeaders", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertNodeInfoHeaders
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertNodeInfoHeaders(conn)
	})

	It("should insert new headers when they don't exist", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Clean up any existing records for this tenant
		err = conn.DB.WithContext(ctx).
			Table("node_info_headers").
			Where("tenant_id = ?", tenant.ID).
			Delete(&table.NodeInfoHeaders{}).Error
		Expect(err).ToNot(HaveOccurred())

		headers := domain.Headers{
			Commerce: "test-commerce",
			Consumer: "test-consumer",
			Channel:  "test-channel",
		}

		err = upsert(ctx, headers)
		Expect(err).ToNot(HaveOccurred())

		// Verify headers were saved
		var savedHeaders table.NodeInfoHeaders
		err = conn.DB.WithContext(ctx).
			Table("node_info_headers").
			Where("document_id = ?", headers.DocID(ctx)).
			First(&savedHeaders).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(savedHeaders.Commerce).To(Equal("test-commerce"))
		Expect(savedHeaders.Consumer).To(Equal("test-consumer"))
		Expect(savedHeaders.Channel).To(Equal("test-channel"))
		Expect(savedHeaders.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should not insert headers if they already exist", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Clean up any existing records for this tenant
		err = conn.DB.WithContext(ctx).
			Table("node_info_headers").
			Where("tenant_id = ?", tenant.ID).
			Delete(&table.NodeInfoHeaders{}).Error
		Expect(err).ToNot(HaveOccurred())

		headers := domain.Headers{
			Commerce: "test-commerce",
			Consumer: "test-consumer",
			Channel:  "test-channel",
		}

		// First insert
		err = upsert(ctx, headers)
		Expect(err).ToNot(HaveOccurred())

		// Try to insert again
		err = upsert(ctx, headers)
		Expect(err).ToNot(HaveOccurred())

		// Verify only one record exists
		var count int64
		err = conn.DB.WithContext(ctx).
			Table("node_info_headers").
			Where("document_id = ?", headers.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))

		// Verify the record has the correct tenant ID
		var savedHeaders table.NodeInfoHeaders
		err = conn.DB.WithContext(ctx).
			Table("node_info_headers").
			Where("document_id = ?", headers.DocID(ctx)).
			First(&savedHeaders).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(savedHeaders.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should handle empty headers", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Clean up any existing records for this tenant
		err = conn.DB.WithContext(ctx).
			Table("node_info_headers").
			Where("tenant_id = ?", tenant.ID).
			Delete(&table.NodeInfoHeaders{}).Error
		Expect(err).ToNot(HaveOccurred())

		headers := domain.Headers{}

		err = upsert(ctx, headers)
		Expect(err).ToNot(HaveOccurred())

		// Verify empty headers were saved
		var savedHeaders table.NodeInfoHeaders
		err = conn.DB.WithContext(ctx).
			Table("node_info_headers").
			Where("document_id = ?", headers.DocID(ctx)).
			First(&savedHeaders).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(savedHeaders.Commerce).To(BeEmpty())
		Expect(savedHeaders.Consumer).To(BeEmpty())
		Expect(savedHeaders.Channel).To(BeEmpty())
		Expect(savedHeaders.TenantID.String()).To(Equal(tenant.ID.String()))
	})
})
