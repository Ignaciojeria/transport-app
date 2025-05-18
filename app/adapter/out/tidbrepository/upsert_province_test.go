package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertProvince", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertProvince
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertProvince(conn)
	})

	It("should insert province if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		province := domain.Province("Santiago")

		err = upsert(ctx, province)
		Expect(err).ToNot(HaveOccurred())

		var dbProvince table.Province
		err = conn.DB.WithContext(ctx).
			Table("provinces").
			Where("document_id = ?", province.DocID(ctx)).
			First(&dbProvince).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbProvince.Name).To(Equal("Santiago"))
		Expect(dbProvince.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update province if fields are different", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.Province("Nombre Original")

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Province("Nombre Modificado")

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbProvince table.Province
		err = conn.DB.WithContext(ctx).
			Table("provinces").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbProvince).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbProvince.Name).To(Equal("Nombre Modificado"))
		Expect(dbProvince.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should not update if no fields changed", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		province := domain.Province("Sin Cambios")

		err = upsert(ctx, province)
		Expect(err).ToNot(HaveOccurred())

		// Get initial timestamp
		var initialProvince table.Province
		err = conn.DB.WithContext(ctx).
			Table("provinces").
			Where("document_id = ?", province.DocID(ctx)).
			First(&initialProvince).Error
		Expect(err).ToNot(HaveOccurred())
		initialUpdatedAt := initialProvince.UpdatedAt

		// Try to update with same values
		err = upsert(ctx, province)
		Expect(err).ToNot(HaveOccurred())

		var dbProvince table.Province
		err = conn.DB.WithContext(ctx).
			Table("provinces").
			Where("document_id = ?", province.DocID(ctx)).
			First(&dbProvince).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbProvince.UpdatedAt).To(Equal(initialUpdatedAt))
		Expect(dbProvince.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should allow same province for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		province1 := domain.Province("Multi Org Province")
		province2 := domain.Province("Multi Org Province")

		err = upsert(ctx1, province1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, province2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := province1.DocID(ctx1)
		docID2 := province2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each province belongs to its respective tenant using DocID
		var dbProvince1, dbProvince2 table.Province
		err = conn.DB.WithContext(ctx1).
			Table("provinces").
			Where("document_id = ?", docID1).
			First(&dbProvince1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbProvince1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("provinces").
			Where("document_id = ?", docID2).
			First(&dbProvince2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbProvince2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no provinces table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		province := domain.Province("Error Case")

		upsert := NewUpsertProvince(noTablesContainerConnection)
		err = upsert(ctx, province)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("provinces"))
	})
})
