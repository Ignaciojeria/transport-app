package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertContact", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertContact
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertContact(conn, nil)
	})

	It("should insert contact if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		contact := domain.Contact{
			FullName:     "Juan Pérez",
			PrimaryEmail: "juan@example.com",
			PrimaryPhone: "123456789",
		}

		err = upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = conn.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", contact.DocID(ctx)).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.FullName).To(Equal("Juan Pérez"))
		Expect(dbContact.Phone).To(Equal("123456789"))
		Expect(dbContact.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update contact if fields are different", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.Contact{
			FullName:     "Nombre Original",
			PrimaryEmail: "update@example.com",
			PrimaryPhone: "111111111",
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Contact{
			FullName:     "Nombre Modificado",
			PrimaryEmail: "update@example.com",
			PrimaryPhone: "222222222",
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = conn.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.FullName).To(Equal("Nombre Modificado"))
		Expect(dbContact.Phone).To(Equal("222222222"))
		Expect(dbContact.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should not update if no fields changed", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		contact := domain.Contact{
			FullName:     "Sin Cambios",
			PrimaryEmail: "same@example.com",
			PrimaryPhone: "333333333",
		}

		err = upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		// Capturar CreatedAt original para comparar después
		var originalRecord table.Contact
		err = conn.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", contact.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Ejecutar de nuevo sin cambios
		err = upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = conn.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", contact.DocID(ctx)).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.FullName).To(Equal("Sin Cambios"))
		Expect(dbContact.Phone).To(Equal("333333333"))
		Expect(dbContact.CreatedAt).To(Equal(originalRecord.CreatedAt)) // Verificar que no se modificó
		Expect(dbContact.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should allow same contact info for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		contact1 := domain.Contact{
			FullName:     "Multi Org",
			PrimaryEmail: "multi@example.com",
			PrimaryPhone: "444444444",
		}

		contact2 := domain.Contact{
			FullName:     "Multi Org",
			PrimaryEmail: "multi@example.com",
			PrimaryPhone: "444444444",
		}

		err = upsert(ctx1, contact1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, contact2)
		Expect(err).ToNot(HaveOccurred())

		// Verify they have different document IDs and belong to different tenants
		Expect(contact1.DocID(ctx1)).ToNot(Equal(contact2.DocID(ctx2)))

		// Verify each contact belongs to its respective tenant
		var dbContact1, dbContact2 table.Contact
		err = conn.DB.WithContext(ctx1).
			Table("contacts").
			Where("document_id = ?", contact1.DocID(ctx1)).
			First(&dbContact1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("contacts").
			Where("document_id = ?", contact2.DocID(ctx2)).
			First(&dbContact2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should insert using NationalID if email and phone are missing", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		contact := domain.Contact{
			FullName:   "Identificado por RUN",
			NationalID: "12345678-9",
		}

		err = upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = conn.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", contact.DocID(ctx)).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.NationalID).To(Equal("12345678-9"))
		Expect(dbContact.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should generate new ReferenceID if all keys are missing", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		contact := domain.Contact{
			FullName: "Sin Identificación",
		}

		err = upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = conn.DB.WithContext(ctx).
			Table("contacts").
			Where("full_name = ?", contact.FullName).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.DocumentID).ToNot(BeEmpty())
		Expect(dbContact.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should fail if database has no contacts table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		contact := domain.Contact{
			FullName:     "Error Esperado",
			PrimaryEmail: "fail@example.com",
		}

		upsert := NewUpsertContact(noTablesContainerConnection, nil)
		err = upsert(ctx, contact)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("contacts"))
	})
})
