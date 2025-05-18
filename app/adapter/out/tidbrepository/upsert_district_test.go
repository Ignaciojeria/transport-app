package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertDistrict", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should insert a new district", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		upsert := NewUpsertDistrict(conn)
		testDistrict := domain.District("Test District")

		err = upsert(ctx, testDistrict)
		Expect(err).To(BeNil())

		var savedDistrict table.District
		err = conn.WithContext(ctx).
			Table("districts").
			Where("document_id = ?", testDistrict.DocID(ctx).String()).
			First(&savedDistrict).Error
		Expect(err).To(BeNil())
		Expect(savedDistrict.Name).To(Equal(testDistrict.String()))
		Expect(savedDistrict.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should create a new record when district name changes", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		upsert := NewUpsertDistrict(conn)
		testDistrict := domain.District("Test District")

		// First insert
		err = upsert(ctx, testDistrict)
		Expect(err).To(BeNil())

		var firstDistrict table.District
		err = conn.WithContext(ctx).
			Table("districts").
			Where("document_id = ?", testDistrict.DocID(ctx).String()).
			First(&firstDistrict).Error
		Expect(err).To(BeNil())
		firstID := firstDistrict.ID
		Expect(firstDistrict.TenantID.String()).To(Equal(tenant.ID.String()))

		// Insert with different name
		updatedDistrict := domain.District("Updated District")
		err = upsert(ctx, updatedDistrict)
		Expect(err).To(BeNil())

		var secondDistrict table.District
		err = conn.WithContext(ctx).
			Table("districts").
			Where("document_id = ?", updatedDistrict.DocID(ctx).String()).
			First(&secondDistrict).Error
		Expect(err).To(BeNil())
		Expect(secondDistrict.ID).ToNot(Equal(firstID)) // Should be a new record
		Expect(secondDistrict.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should handle multiple districts with different DocIDs", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		upsert := NewUpsertDistrict(conn)
		district1 := domain.District("District 1")
		district2 := domain.District("District 2")

		err = upsert(ctx, district1)
		Expect(err).To(BeNil())
		err = upsert(ctx, district2)
		Expect(err).To(BeNil())

		var count int64
		err = conn.WithContext(ctx).
			Table("districts").
			Where("tenant_id = ?", tenant.ID).
			Count(&count).Error
		Expect(err).To(BeNil())
		Expect(count).To(Equal(int64(2)))
	})

	It("should handle database errors gracefully", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		upsert := NewUpsertDistrict(conn)
		// Create an invalid district with a very long name to trigger a database error
		invalidDistrict := domain.District("This is a very long district name that exceeds the maximum length allowed by the database schema and should cause an error when trying to insert it into the database tableThis is a very long district name that exceeds the maximum length allowed by the database schema and should cause an error when trying to insert it into the database table")

		err = upsert(ctx, invalidDistrict)
		Expect(err).NotTo(BeNil())
	})

	It("should insert and retrieve an empty record", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		upsert := NewUpsertDistrict(conn)
		// Insert empty district
		emptyDistrict := domain.District("")
		err = upsert(ctx, emptyDistrict)
		Expect(err).To(BeNil())

		// Try to retrieve the empty record
		var savedDistrict table.District
		err = conn.WithContext(ctx).
			Table("districts").
			Where("document_id = ?", emptyDistrict.DocID(ctx).String()).
			First(&savedDistrict).Error
		Expect(err).To(BeNil())
		Expect(savedDistrict.Name).To(Equal(""))
		Expect(savedDistrict.DocumentID).To(Equal(emptyDistrict.DocID(ctx).String()))
		Expect(savedDistrict.TenantID.String()).To(Equal(tenant.ID.String()))
	})
})
