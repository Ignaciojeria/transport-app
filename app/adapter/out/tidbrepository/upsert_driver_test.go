package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertDriver", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertDriver
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertDriver(conn, nil)
	})

	It("should insert driver if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		driver := domain.Driver{
			Name:       "Juan Pérez",
			NationalID: "12345678-9",
		}

		err = upsert(ctx, driver)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := driver.DocID(ctx)

		// Verify the record was inserted correctly
		var dbDriver table.Driver
		err = conn.DB.WithContext(ctx).
			Table("drivers").
			Where("document_id = ?", docID).
			First(&dbDriver).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbDriver.Name).To(Equal("Juan Pérez"))
		Expect(dbDriver.NationalID).To(Equal("12345678-9"))
		Expect(dbDriver.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update driver if fields are different", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.Driver{
			Name:       "Nombre Original",
			NationalID: "11111111-1",
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := original.DocID(ctx)

		modified := domain.Driver{
			Name:       "Nombre Modificado",
			NationalID: "11111111-1", // Mismo NationalID
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		// Verify the record was updated correctly
		var dbDriver table.Driver
		err = conn.DB.WithContext(ctx).
			Table("drivers").
			Where("document_id = ?", docID).
			First(&dbDriver).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbDriver.Name).To(Equal("Nombre Modificado"))
		Expect(dbDriver.NationalID).To(Equal("11111111-1"))
		Expect(dbDriver.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should not update if no fields changed", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		driver := domain.Driver{
			Name:       "Sin Cambios",
			NationalID: "22222222-2",
		}

		err = upsert(ctx, driver)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocID
		docID := driver.DocID(ctx)

		// Capturar CreatedAt original para comparar después
		var originalRecord table.Driver
		err = conn.DB.WithContext(ctx).
			Table("drivers").
			Where("document_id = ?", docID).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Ejecutar de nuevo sin cambios
		err = upsert(ctx, driver)
		Expect(err).ToNot(HaveOccurred())

		// Verify the record was not updated
		var dbDriver table.Driver
		err = conn.DB.WithContext(ctx).
			Table("drivers").
			Where("document_id = ?", docID).
			First(&dbDriver).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbDriver.Name).To(Equal("Sin Cambios"))
		Expect(dbDriver.CreatedAt).To(Equal(originalRecord.CreatedAt)) // Verificar que no se modificó
		Expect(dbDriver.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should allow same driver for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		driver1 := domain.Driver{
			Name:       "Multi Org Driver",
			NationalID: "33333333-3",
		}

		driver2 := domain.Driver{
			Name:       "Multi Org Driver",
			NationalID: "33333333-3",
		}

		err = upsert(ctx1, driver1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, driver2)
		Expect(err).ToNot(HaveOccurred())

		// Get the DocIDs
		docID1 := driver1.DocID(ctx1)
		docID2 := driver2.DocID(ctx2)

		// Verify they have different document IDs
		Expect(docID1).ToNot(Equal(docID2))

		// Verify each driver belongs to its respective tenant using DocID
		var dbDriver1, dbDriver2 table.Driver
		err = conn.DB.WithContext(ctx1).
			Table("drivers").
			Where("document_id = ?", docID1).
			First(&dbDriver1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbDriver1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("drivers").
			Where("document_id = ?", docID2).
			First(&dbDriver2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbDriver2.TenantID.String()).To(Equal(tenant2.ID.String()))

		// Verify total count of records with same name and national ID
		var count int64
		err = conn.DB.WithContext(context.Background()).
			Table("drivers").
			Where("document_id IN (?, ?)", docID1, docID2).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	It("should fail if database has no drivers table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		driver := domain.Driver{
			Name:       "Error Esperado",
			NationalID: "44444444-4",
		}

		upsert := NewUpsertDriver(noTablesContainerConnection, nil)
		err = upsert(ctx, driver)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("drivers"))
	})
})
